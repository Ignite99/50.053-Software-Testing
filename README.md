# HTTP, CoAP, BLE Fuzzer (Team 5)

## Group Members

| Member Name  | Student ID | Year |
| ------------- | ------------- | ------------- |
| Andrean Priadi | 1005632  | Senior |
| Eliana Setiabudi  | 1005252 | Senior |
| Nathan Chang | 1005149  | Senior |
| Nicholas Goh | 1005194 | Senior |
| Nur Thohirah Bte Sani | 1005463 | Senior |

## About

### Fuzzer Overview
This project aims to create a fuzzer for the black-box testing of three target applications utilising Django (HTTP), CoAP, and BLE. To achieve this goal, the team delved into HTTP, CoAP, and BLE message formats, researched various black-box fuzzing techniques, and successfully developed a fuzzer written in Go programming language to test these applications.

### Fuzzer Outcomes
The following are the outcomes of our fuzzer:
- Ease of user input and fuzzer usage
- Tunable parameters such as energy assignment
- Fuzz HTTP, CoAP, BLE payloads
- Send HTTP, CoAP, BLE requests to servers
- Receive fuzzed payload responses from servers
- Record server requests and responses
- Generalisation of fuzzing implementation for various targets
- Categorise server responses based valid and invalid conventions
- Fuzz a certain input multiple times according to its assigned energy.
- Implementation of IsInteresting
- Comprehensive logger to show fuzzing results


## Running the Fuzzer

### Fuzzing Django / HTTP
1. Run HTTP server 
2. `cd golang_src`
3. Edit `inputs/add_product.json` according to your desired request body.
E.g. 
{
  "name": "chicken",
  "info": "abcd",
  "price": "1000"
}
4. Then run: `go run main.go DJANGO http://127.0.0.1:8000/datatb/product/add/ POST ./inputs/add_product.json 200`
5. The output logger, `log.html`, will be updated in `golang_src/fuzzing_responses`.

### Fuzzing CoAP
1. Assuming you have python2 and CoAPthon dependencies already, `cd` to the CoAPthon directory
2. Run the server with `python2 coapserver.py -i 127.0.0.1 -p 5683`
3. Next, open the fuzzer directory and `cd golang_src`
4. Enter `go run main.go COAP 127.0.0.1:5683 POST ./inputs/hello_world.txt 3`     
5. The output loggers, `log.html` and `unique_logs.html`, will be updated in `golang_src/fuzzing_responses`.

### Fuzzing BLE
1. `cd golang_src`
2. `go run main.go BLE ./ POST ./ 200`
3. `GCOV_PREFIX=$(pwd) GCOV_PREFIX_STRIP=3 <path/to/zephyr.exe> --bt-dev=127.0.0.1:9000`

