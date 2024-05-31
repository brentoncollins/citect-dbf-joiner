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

| Option                                 | Comment                                                                                                                                                                                                | Default               |
|----------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------------------|
| -input-dir                             | The input directory that either contains the master.dbf or multiple compile folders (required)  **(required)**                                                                                         | NA                    |
| -dbf                                   | Comma seperated list of DBFs you want to join, or single string, each one will have a separate output.   **(required)**                                                                                | variable,digalm,equip |
| -output-dir                            | The output directory, if multiple dbf files are specified, a separate csv will be created for each                                                                                                     | Current working dir   |
| -log-path                              | The output log path                                                                                                                                                                                    | Current working dir   |
| -master-dbf-dir                        | Set this flag to use the master DBF file in the input-dir to determine the folders that contain the DBF files.                                                                                         | False                 |
| -search-compile-folders-for-master-dbf | Set this flag to find all sub-folders in input-dir, then search each sub-folder by date descending for a master DBF to determine the folders that contain the DBF files DBF files within that folder . | False                 |
| -project-compile-dir                   | Set this flag to find all sub-folders in input-dir, and get DBFs from each folder.                                                                                                                     | False                 |
| -search-latest-compile-dirs            | Set this flag to find all sub-folders in input-dir, identify the latest sub-folder by date descending, each sub-folder within the latest sub-folder will be searched for a DBF.                        | False                 |


The minimum required is
- -input-dir

Choose one of these, see below of explanation.
- -master-dbf-dir
- -search-compile-folders-for-master-dbf
- -project-compile-dir
- -search-latest-compile-dirs

```
citect-dbf-joiner.exe -input-dir "C:\ProgramData\AVEVA Plant SCADA 2023\User" -search-latest-compile-dirs

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
citect-dbf-joiner.exe -input-dir "C:\ProgramData\AVEVA Plant SCADA 2023\User" -out-putdir "C:\Temp\" -log-path "C:\Temp\application.log" -dbf "variable" -search-latest-compile-dirs

Output
Successfully written to file: C:\\Temp\\variable.csv
```

### Example of input directory structure

#### Master DBF in input-dir flags
- input-dir="C:\Some\Folder\User"
- master-dbf-dir

```
Set this flag to use the master DBF file in the input-dir (User), read the master DBF to determine the folders that contain the DBF files.
.
└── User
        ├── folder1
        │   └── variable.dbf    * Found this file
        ├── folder2
        │   └── variable.dbf    * Found this file
        ├── folder3
        │   └── variable.dbf    * Found this file
        └── master.dbf
```
#### Master DBF in latest subdirectory flags
- input-dir="C:\Some\Folder\Compiles"
- search-compile-folders-for-master-dbf

```
Set this flag to find all sub-folders (Compile1, Compile2) in input-dir (Compiles), identify the latest sub-folder 
by date descending (Compile2), read the master DBF to identify all the sub-folders (folder1, folder2, folder3) 
to find the DBF files.

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
            │ 	└── variable.dbf    * Found this file
            ├── folder2
            │ 	└── variable.dbf    * Found this file
            ├── folder3
            │ 	└── variable.dbf    * Found this file
            └── master.dbf
```

#### No master DBF, use folder of input-dir to find DBF files.
- input-dir="C:\Some\Folder\User"
- project-compile-dir

```
Set this flag to use all subfolders in input-dir (User), loop throught each one to find the DBF files.
.
└── User
        ├── folder1
        │   └── variable.dbf    * Found this file
        ├── folder2
        │   └── variable.dbf    * Found this file
        └── folder3
            └── variable.dbf    * Found this file
```

#### No master DBF, use latest sub-folder of input-dir to find subdirectories with DBF files.
- input-dir="C:\Some\Folder\Compiles"
- search-compile-folders-for-master-dbf

```
Set this flag to find all sub-folders (Compile1, Compile2) in input-dir (Compiles), identify the latest sub-folder 
by date descending (Compile2), loop through all subfolders within the latest sub-folder (folder1, folder2, folder3) 
to find the DBF files.

.
└── Compiles
        │
        ├── Compile1------------------------------Date: 01/01/2021
        │   ├── folder1
        │   │   └── variable.dbf
        │   ├── folder2
        │   │   └── variable.dbf
        │   └── folder3
        │      └── variable.dbf
        │   
        │   
        └── Compile2------------------------------Date: 02/01/2021
            ├── folder1
            │ 	└── variable.dbf    * Found this file
            ├── folder2
            │ 	└── variable.dbf    * Found this file
            └── folder3
             	└── variable.dbf    * Found this file

```

### Important Note

This project uses the `go-dbase/dbase` library for handling DBF files.
Please note that this library does not officially support DBF version 3, it supports version 5, which is used by Citect.
However, reading operations should work as expected.
