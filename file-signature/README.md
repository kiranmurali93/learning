# File Signatures

Generally files have few set of numbers called magic numbers. These numbers help us in determining the actual file type and how to treat it.

## Project Structure
| - main.go
| - signatures.go

## Running the demo
1. ```go run . <filepath/fileName>```

## File purposes

### ```main.go```

This file has the logic for taking out the headers of the file and comparing with the standard magic numbers associated with different file.

### ```signatures.go```

This file has the list of supported file type and their magic numbers defined.

Along with this it also has a offset for which the header check has to be done.

### ```scripts_check.go```

This file is an example for if we want to add additional logic for checking the file signature