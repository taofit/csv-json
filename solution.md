## Solution

When clients send/post data in text/csv format, firstly the backend server will save all the records in a two dimensinal slice. Values in each row of the cvs file is saved in a slice that is also the element of the outer slice.

Then the programm will loop each row and check if the file header is in the right order and if elements in each row follow the right structure. If all the rows are valid, each value in every row will be saved in a slice of struct type. the struct contains element value and its parent path. The parent path is the concatenation of each value of element comes before the current element in a row. Such as each element in level_1 column has parent path of empty string as there is no element coming before it, element in level_2 column has parent path whose value is level_1, and so on.

Then programm will start build a Node struct by checking all the elements in the record. It will first look for the elements without parent path, and same them as the first level children in the Node, then it will look for the elements whose parents are the first level children in the Node, and so on and so forth until reaching all the elements who have no children. It is achieved by a recursive call. The way to find children of each element is to compare the parent path of each element in the record with current element's path(its own parent path + its value)

The programm build a node with multi level children, it is not just limited to 3 levels, could be 4 or 5 level. It is mainly achieved by calling the recursively function when building the node.

## Run and test it

To run the programm, first go to the cmd folder and run command: `go run main.go` which is to start the backend http server, then send post API request with csv formatted plain text via postman or insomnia to the server http://localhost:8080/

Also the programm provides test file, and it can be run via command: `go test -v` in the root folder
