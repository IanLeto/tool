package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func generateSpans(ctx context.Context, tracer trace.Tracer, depth int) {
	if depth == 0 {
		return
	}

	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	_, span := tracer.Start(ctx, fmt.Sprintf("span-%d", depth))
	defer span.End()

	// 模拟服务处理时间
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

	// 添加一些属性
	span.SetAttributes(
		attribute.String("service", fmt.Sprintf("service-%d", depth)),
		attribute.Int("depth", depth),
	)

	// 随机决定是否产生错误
	if rand.Intn(5) == 0 {
		err := fmt.Errorf("random error at depth %d", depth)
		span.RecordError(err)
		cancel(err)
		return
	}

	generateSpans(ctx, tracer, depth-1)
}

var SpanCmd = &cobra.Command{
	Use: "span",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			file   *os.File
			err    error
			format string
		)

		path, _ := cmd.Flags().GetString("path")
		format, _ = cmd.Flags().GetString("format")

		if path == "" {
			file = os.Stdout
		} else {
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0755)
				NoErr(err)
			}
			file, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		}

		var exporter sdktrace.SpanExporter
		switch format {
		case "compact":
			exporter, err = stdouttrace.New(
				stdouttrace.WithWriter(file),
				stdouttrace.WithPrettyPrint(),
				stdouttrace.WithoutTimestamps(),
			)
		case "escape":
			exporter, err = newEscapeExporter(file)
		default:
			exporter, err = stdouttrace.New(
				stdouttrace.WithWriter(file),
				stdouttrace.WithPrettyPrint(),
			)
		}
		if err != nil {
			fmt.Println("failed to create exporter:", err)
			os.Exit(1)
		}

		// 创建 tracer provider
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(resource.NewWithAttributes(
				"goOri", attribute.String("service.name", "trace-generator"),
			)),
		)
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				err = errors.Join(err, fmt.Errorf("failed to shutdown TracerProvider"))
				fmt.Println(err)
			}
		}()
		otel.SetTracerProvider(tp)

		// 生成跟踪数据
		tracer := otel.Tracer("trace-generator")
		for i := 0; i < 10; i++ {
			ctx, cancel := context.WithCancelCause(context.Background())
			_, span := tracer.Start(ctx, fmt.Sprintf("trace-%d", i))
			generateSpans(ctx, tracer, rand.Intn(5)+1)
			span.End()
			cancel(nil)
		}
	},
}

func init() {
	SpanCmd.Flags().StringP("path", "p", "", "path")
	SpanCmd.Flags().StringP("format", "f", "default", "output format, options: default|compact|escape")
}

type escapeExporter struct {
	w io.Writer
}

func newEscapeExporter(w io.Writer) (sdktrace.SpanExporter, error) {
	return &escapeExporter{w: w}, nil
}

func (e *escapeExporter) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	for _, span := range spans {
		_, err := e.w.Write([]byte(strings.ReplaceAll(fmt.Sprintf("%#v\n", span), " ", "")))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *escapeExporter) Shutdown(ctx context.Context) error {
	return nil
}
