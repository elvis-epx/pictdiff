use std::env;
use std::process;
extern crate image;
use image::GenericImage;
use image::ImageBuffer;
use image::Rgba;
use std::cmp;
use std::path::Path;

fn pre_mult_alpha(oldp: &Rgba<u8>, newp: &Rgba<u8>, channel: usize) -> i64
{
	let old = (oldp[channel] as f64) * ((oldp[3] as f64) / 255.0);
	let new = (newp[channel] as f64) * ((newp[3] as f64) / 255.0);
	(new as i64) - (old as i64)
}

fn main()
{
	let minute: i64 = 5;
	let increase_minute: i64 = 2;
	let mut totaldiff: u64 = 0;
	
	let args: Vec<_> = env::args().collect();
	if args.len() < 4 {
		println!("Usage: {} oldpicture newpicture diffmap", &args[0]);
		process::exit(2);
	}

	let img1 = image::open(&Path::new(&args[1])).unwrap();
	let img2 = image::open(&Path::new(&args[2])).unwrap();
	
	if img1.dimensions() != img2.dimensions() {
		println!("Images have different sizes, cannot compare");
		process::exit(1);
	}
	
	let (width, height) = img1.dimensions();

	let img1 = img1.to_rgba();
	let img2 = img2.to_rgba();
	
	let mut imgmap = ImageBuffer::new(width, height);
	
	for (x, y, mappixel) in imgmap.enumerate_pixels_mut() {
       		let p1 = img1.get_pixel(x, y);
       		let p2 = img2.get_pixel(x, y);
		let mut diffpixel: [i64; 3] = [255, 255, 255];
	
		let diffs: [i64; 4] = [pre_mult_alpha(&p1, &p2, 0),
					pre_mult_alpha(&p1, &p2, 1),
					pre_mult_alpha(&p1, &p2, 2),
					(p2[3] as i64) - (p1[3] as i64)];

		let absdiff = diffs[0].abs() + diffs[1].abs() + diffs[2].abs() + diffs[3].abs();
		totaldiff += absdiff as u64;
		let diffplus: [i64; 3] = [cmp::max(0, diffs[0]),
			cmp::max(0, diffs[1]), cmp::max(0, diffs[2])];
		let totplus = diffplus[0] + diffplus[1] + diffplus[2];
		let diffminus: [i64; 3] = [cmp::min(0, diffs[0]),
			cmp::min(0, diffs[1]), cmp::min(0, diffs[2])];
	
		for i in 0..3 {
			diffpixel[i] += diffminus[i];
			diffpixel[i] -= totplus;
			diffpixel[i] += diffplus[i];
			if absdiff > 0 && absdiff < minute {
				diffpixel[i] -= increase_minute;
			}
			diffpixel[i] = cmp::max(0, diffpixel[i]);
		}
		let diffpixel8: [u8; 3] = [diffpixel[0] as u8,
						diffpixel[1] as u8,
						diffpixel[2] as u8];
		*mappixel = image::Rgb(diffpixel8);
	}
	
	let _ = imgmap.save(&args[3]).unwrap();
	println!("{}", totaldiff);
}
