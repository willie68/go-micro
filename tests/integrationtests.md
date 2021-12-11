## Common test information for Go services
 (A simple step-to-step guide)    
   
   
 Checkout the project.  
  On your command line switch to the root directory of this project.  


**Unit tests**  
The test creation and execution follows GO convention. [Here is a good introduction into writing tests with Golang.](https://blog.alexellis.io/golang-writing-unit-tests/)  
We recommend to place unit tests into the particular package of the tested go file. In this way code covarage for tests is possible without any problems.  
  Type ->  go test -v ./...     (executes all tests)  
  Type -> go test -v ./... | go-junit-report >report.xml (executes all tests and generates a report)  
  Type -> go test -v -cover ./... | go-junit-report >report.xml (executes all tests with covarage infos and generates a report)  
  Type -> go test -v ./... -coverprofile coverage.txt (executes all tests and creates a file with coverage informations)  
  Type -> go tool cover -html=coverage.txt (opens a HTML page based on a created coverage.txt file with covered and non-covered code marked in terms of color) 

**Integration tests**  
We do not analyze code coverage in our integration tests. So we suggest to place them in a separate folder. (see our example in tests/integrationtests/..)  
As integrationtests require a running system,  they are skipped by default to not interfere with a quick local execution of all unit tests.  
To execute them you have to add **-test.skip=false** as a commandline parameter.  
By default all REST-calls are executed against the proactcloud system. If you want to execute them on your computer you have to provide access to the proact.cloud.  
If you want to make calls to a locally running system add **-test.local=true** as commandline parameter.  
**Please note:** the tests then require the service to be available under https://127.0.0.1:8443.   
Type: go test -v ./tests/integrationtests -test.skip=false (executes only the integrationtests against proact.cloud system)  
Type: go test -v ./tests/integrationtests -test.skip=false -test.local=true | go-junit-report>report.xml(executes only the integrationtests against your local system and generates a report)
