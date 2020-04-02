/*
rasterization code extracted and converted to go from https://github.com/petrkutalek/png2pos/blob/master/png2pos.c
*/

package escpos

const (
	maxY = 1662
)

// Raster rasterize image
func (e *Escpos) Raster(width, height, bytesWidth int, imgBw []byte) {
	flushCmd := []byte{
		/* GS ( L, Print the graphics data in the print buffer,
		   p. 241 Moves print position to the left side of the
		   print area after printing of graphics data is
		   completed */
		0x1d, 0x28, 0x4c, 0x02, 0x00, 0x30,
		/* Fn 50 */
		0x32,
	}

	for l := 0; l < height; {
		nLines := maxY
		if nLines > height-l {
			nLines = height - l
		}

		f112P := 10 + nLines*bytesWidth
		storeCmd := []byte{
			/* GS 8 L, Store the graphics data in the print buffer
			   (raster format), p. 252 */
			0x1d, 0x38, 0x4c,
			/* p1 p2 p3 p4 */
			byte(f112P), byte(f112P >> 8),
			byte(f112P >> 16), byte(f112P >> 24),
			/* Function 112 */
			0x30, 0x70, 0x30,
			/* bx by, zoom */
			0x01, 0x01,
			/* c, single-color printing model */
			0x31,
			/* xl, xh, number of dots in the horizontal direction */
			byte(width), byte(width >> 8),
			/* yl, yh, number of dots in the vertical direction */
			byte(nLines), byte(nLines >> 8),
		}

		e.WriteRaw(storeCmd)
		e.WriteRaw(imgBw[l*bytesWidth : (l+nLines)*bytesWidth])
		e.WriteRaw(flushCmd)

		l += nLines
	}
}
