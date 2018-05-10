# AdminSorter
This tool is specifically designed to automate scanning large volumes of PDF's to a pre-organized directory tree.
Based on the barcode in the image, the filename is automatically named according to the below convention.  If the subfolder does not
exist for this range, then it's automatically created.

## Details
A pre-defined folder with subfolders using the syntax:

Subfolder names:
```
100-500
501-1000
1001-1500
```
Scanned filenames with 1.pdf, 2.pdf, etc... automatically get moved to the subfolder "100-500".
Scanned filenames with 501.pdf, 502.pdf, etc... automatically get filed in the subfolder "501-1000".

