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

Alpha channel of images is not taken into consideration yet. Not sure
how to represent changes in alpha on the diff image.

Differences in 48-bit images may be exagerated in diff image. Images with
different color dephts are not correctly compared.

# Contact

Elvis Pfützenreuter - epxx@epxx.co - https://epxx.co
