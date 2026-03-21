## Filters

Filters are light-weight transformations applied after stages, often to add visual effects.

### Available filters
- `filters/drawedge`: Combines a base ASCII stream with edge detection results.
- `filters/charcolor`: Applies color data to ASCII chars.

### How to create a custom filter
1. Implement a filter that receives stage output types and returns modified output.
2. In command pipeline, use `flow.Mask` with `flow.Zip` to join streams, e.g.: `flow.Mask(ctx, flow.Zip(ctx, streamA, streamB), myFilter)`.

Example skeleton:
```go
package myfilter

import "github.com/26in26/p02-ascii-generator/pipeline"

func NewMyFilter() pipeline.Filter[*image.AsciiBuffer, *image.RGBBuffer] {
    return pipeline.FilterFunc(func(a *image.AsciiBuffer, b *image.RGBBuffer) (*image.AsciiBuffer, error) {
        // merge or modify a with b
        return a, nil
    })
}
```

### Put it together
In `command.go`, filters are attached after stage outputs:
1. Perform common stages (resize, grayscale, ascii)
2. Use `flow.Zip` for combined inputs
3. Use `flow.Mask` to apply filter

Filters are ideal for post-processing that needs multiple input streams.