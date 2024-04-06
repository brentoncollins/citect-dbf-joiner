# _Citect DBF Joiner_

This project is a Go-based utility for working with Citect SCADA (Aveva Plant SCADA) DBF files and output as a CSV. It is designed to find and join DBF files 
from a given compile.

## Compilation

This project is written in Go. To compile it, you'll need to have Go installed on your system. You can download Go from the [official website](https://golang.org/dl/).
Once you have Go installed.

```bash
git clone https://github.com/brentoncollins/citect-dbf-joiner.git
cd citect-dbf-joiner
go build -o citect-dbf-joiner.exe cmd/citect-dbf-joiner/main.go
```

## Usage

This is a command-line utility that can be used to join DBF files from a Citect SCADA compile.

### Command-line reference

| Option     | Comment                                                                                                                  | Default               |
|------------|--------------------------------------------------------------------------------------------------------------------------|-----------------------|
| -inputdir  | The input directory that either contains the master.dbf or multiple compile folders (required)  **(required)**           | NA                    |
| -dbf       | Comma seperated list of DBF's you want to join, or single string, each one will have a seperate output.   **(required)** | variable,digalm,equip |
| -outputdir | The output directory, if multiple dbf files are specified, a separate csv will be created for each                       | Current working dir   |
| -logpath   | The output log path                                                                                                      | Current working dir   |


The minimum required is -inputdir where the DBF files are located.
```
citect-dbf-joiner.exe -inputdir "C:\ProgramData\AVEVA Plant SCADA 2023\User"

Output
Logs...
Successfully written to file: C:\\Temp\\variable.csv
Logs...
Successfully written to file: C:\\Temp\\digalm.csv
Logs...
Successfully written to file: C:\\Temp\\equip.csv
```

Example with all arguments
```
citect-dbf-joiner.exe -inputdir "C:\ProgramData\AVEVA Plant SCADA 2023\User" -outputdir "C:\Temp\" -logpath "C:\Temp\application.log" -dbf "variable"

Output
Successfully written to file: C:\\Temp\\variable.csv
```

##### Example of input directory structure
```
The user can pass the input directory that contains the master.dbf file in the root directory.
.
└── User
        ├── folder1
        │   └── variable.dbf
        ├── folder2
        │   └── variable.dbf
        ├── folder3
        │   └── variable.dbf
        └── master.dbf

Or user can pass the folder with the latest compiles, if the utility does not find the master.dbf file in the root directory mentioned above.
It will attempt to find it in the latest subdirectory based on modification time. 
It will only join in the DBF's in the subfolders of the first found master.dbf file.
.
└── Compiles
        │
        ├── Compile1------------------------------Date: 01/01/2021
        │   ├── folder1
        │   │   └── variable.dbf
        │   ├── folder2
        │   │   └── variable.dbf
        │   ├── folder3
        │   │   └── variable.dbf
        │   └── master.dbf
        │   
        └── Compile2------------------------------Date: 02/01/2021
            ├── folder1
            │ 	└── variable.dbf
            ├── folder2
            │ 	└── variable.dbf
            ├── folder3
            │ 	└── variable.dbf
            └── master.dbf
```

### Important Note

This project uses the `go-dbase/dbase` library for handling DBF files.
Please note that this library does not officially support DBF version 3, it supports version 5, which is used by Citect.
However, reading operations should work as expected.

### Disclaimer

I am no Citect expert, I just use this data for a project I am working on. If you have any suggestions or improvements, please feel free to create a pull request.