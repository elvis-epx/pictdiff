# pictdiff

This program implements image comparison. It prints a simple metric
of difference, and generates an image for visual inspection.

There is a Python implementation, that is our reference implementation;
and a Rust implementation that is 40x faster.

How to use:

```
./pictdiff.py img1.png img2.png diff.png
```

The diff image is white, with dark spots where differences are found.
The spots are darker as the differences are more intense. The spots
are also tinted to reflect the color changes. For example, if "img1"
is generally more red than "img2", then "diff" will be cyan. On the
other hand, if "img2" is more red, then "diff" is also red.

Differences in alpha (opacity) are counted in two ways. They are added
to the total absolute difference, which are rendered as gray in diff
image. This guarantees that alpha differences will be revealed clearly.

Also, the color channels are multiplied by alpha, so the color differences
are mitigated when both images are quite transparent. For example, two
completely transparent images (alpha=0) will be considered equal, even
if one is "red" and another is "blue".

How to run the Rust version:

```
cargo run --release img1.png img2.png diff.png
```

# Motivation

The 'compare' tool from ImageMagick almost does what I want, but it
marks any difference as a red pixel. The threshold can be configured,
but the generated diff image is an all-or-nothing comparison. I needed
a diff image for quick inspection of changes, and the image should show
the magnitude of these differences, as well as how color was affected.

# Bugs

Not tested with 48-bit images. (Python Imaging Library does not open
48-bit TIFFs; Rust image crate does not support 48-bit yet.)

# Contact

Elvis Pfützenreuter - epxx@epxx.co - https://epxx.co
