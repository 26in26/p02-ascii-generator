## Stages

Stages are building blocks for image transformations. Each stage implements the generic pipeline stage interface and is used through `flow.Map`.

### Available stages
- `stages/resize`: Resizes the image. Configurable via `resize.With`, `resize.WithWidth`, `resize.WithHeight`, `resize.WithAspectRatio`.
- `stages/grayscale`: Converts RGB to grayscale buffer.
- `stages/ascii`: Converts grayscale buffer to `AsciiBuffer` using charset density and invert options.
- `stages/edge`: Detects edges using Sobel filter.

### How to create a custom stage
1. Implement a stage type satisfying `pipeline.Stage[input, output]`.
2. Provide options and constructors like `New<Stage>Name`.
3. Use in your code via `flow.Map(ctx, &inputOutlet, yourStage)`.

Example skeleton:
```go
package mystage

import (
    "context"
    "github.com/26in26/p02-ascii-generator/image"
    "github.com/26in26/p02-ascii-generator/pipeline"
)

func NewMyStage() pipeline.Stage[*image.RGBBuffer, *image.RGBBuffer] {
    return &struct myStage{
        a: 1,
        b: 2,
    } 
}

type myStage struct { // staisfies the stage interface
    a int
    b int
}

func (s *myStage) Process(ctx context.Context, input *image.RGBBuffer) (*image.RGBBuffer, error) {...}
```

### Usage in pipeline
`flow.Map(ctx, &resizeStream, customStage)`

### Note
- Stages can release their output alocated memory when it's not needed through the `Release()` method.
- Stage should return errors on invalid input, the pipeline handles propagation.
