# P02 — Ascii Generator

> 26 in 26 · Weeks 03–04 · image processing
<p align="center"><img width="450" height="300" alt="worship rust" src="assets/ascii-generator.png" /></p>

## Goal
Build a high-performance package that encapsulates image processing techniques to generate ASCII art from images with configurable pipelines and algorithms.

## Scope
**In scope**
- Core image-to-ASCII conversion engine
- Configurable processing pipeline (sampling, resizing, dithering, mapping)
- Multiple character set options (density levels, custom sets)
- Performance optimization for real-time processing
- Support for common image formats

**Out of scope**
- GUI/frontend implementation
- Video streaming processing
- Advanced color-to-grayscale algorithms beyond standard methods

## Timeline
- **Week 1:** Design, research, POC
- **Week 2:** Implementation, testing, docs

## Status
- [X] Design
- [X] POC
- [ ] Core implementation
- [ ] Tests
- [ ] Documentation

## 🛠 Tech Stack
- Language: Go
- Constraints: no external image processing packages

## 🚀 Running the Project
```bash
# example
make run
```

## 📋 API Overview
- Image loading and preprocessing
- Configurable algorithm pipeline
- Character mapping strategies
- Output formatting options
