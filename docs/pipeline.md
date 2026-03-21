## Pipeline Overview

A pipeline is the core processing chain that takes an input image and transforms it into ASCII art. In this library, a pipeline is built from:
- `Outlets` (data flow channels)
- `Stages` (image processing transformations)
- `Filters` (post-stage modifications e.g. edge drawing, color mapping)

A pipeline is operated as asynchronous stream processing, where each stage processes data in sequence and data flows through channels managed by `pipeline/flow`.

### 1. Outlets
An `Outlet` is a typed stream abstraction. In `command.go`, `rawStream` is the input outlet, and the final transformed stream is consumed by `handleOutput`.

*outlet utilities:*
- flow.map - map is a function that trancent one outlet to another
- flow.zip - zip combines to outlets into one with unified type
- flow.mask - mask is the map equivelent for filters


### 2. Stages
A stage is a processing step. Common stages include:
- resize (from `stages/resize`)
- grayscale (from `stages/grayscale`)
- ascii mapping (from `stages/ascii`)
- edge detection (from `stages/edge`)

Stages can be chained using `flow.Map`.

### 3. Filters
Filters are applied with `flow.Mask` and can combine outputs from multiple stages via `flow.Zip`.
Example in `command.go`:
- `drawedge.NewEdgeFilter` adds edge lines to ASCII output
- `charcolor.NewColorFilter` colors ASCII characters

### Stage vs Filter
- `Stage` transforms image buffers; it is the main pipeline logic.
- `Filter` applies additional modifications (binary mask, color overlay) on already-stage processed output.

### Recommended Flow
1. Load image into `RGBBuffer`
2. Resize stage
3. Grayscale stage
4. ASCII stage
5. Optional edge/color filters
6. Output to terminal or file