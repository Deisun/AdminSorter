# AdminSorter
This tool is designed to monitor and automate organization of large volumes of scanned PDF documents.
Filenames are automatically named using barcode recognition.  They are then sorted in to subfolders according to filename.  If the subfolder does not exist, then it's automatically created.

## Details
Syntax for filenames:
```
1.pdf
2.pdf
3.pdf
...
```

Subfolder names:
```
0-500
501-1000
1001-1500
...
```
Scanned filenames with 1.pdf, 2.pdf, etc... automatically get moved to the subfolder "100-500".
Scanned filenames with 501.pdf, 502.pdf, etc... automatically get filed in the subfolder "501-1000".

