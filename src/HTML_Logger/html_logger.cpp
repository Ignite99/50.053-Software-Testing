#include <iostream>
#include <csignal>
#include <fstream>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>

using namespace std;

class HTMLLogger {
    public:
        string project_type;
        string output_file_path;
        string output_file_name;
        ofstream output_file;
        string header = 
            R"(<!DOCTYPE html>
            <html lang="en" class="scroll-smooth" dir="ltr">
                <head>
                    <meta charset="UTF-8" />
                    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
                    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
                    <link href="output_style.css" rel="stylesheet">
                    <title>Fuzzer Output Log</title>
                </head>
                <body>)";
        string footer = 
            R"(     </table>
                </body>
            </html>
            )";
        string table_columns=
            R"(
                <table>
                    <tr style="background-color:lightgrey;">
                        <th>Time</th>
                        <th>Request Type</th>
                        <th>Contents</th>
                        <th>Output</th>
                    </tr>
            )";

    public:
        HTMLLogger(string output_file_path, string output_file_name, string project_type){
            this->output_file_path = output_file_path;
            this->output_file_name = output_file_name;
            this->project_type = project_type;
        }

        void create_file(){
            output_file.open(output_file_name);
            if(!output_file){
                cout << "Creating file " << output_file_name << " ..." << endl;
                output_file.open(output_file_name);
            }
            output_file << header;
            if(project_type == "DJANGO"){
                output_file << "<p><b>Fuzz Target: </b>Django</p>";
            } else if(project_type == "COAP"){
                output_file << "<p><b>Fuzz Target: </b>CoAP</p>";
            } else if(project_type == "BLE"){
                output_file << "<p><b>Fuzz Target: </b>BLE</p>";
            }
            // TODO - display any other overall stats
            output_file << table_columns;
        }

        void add_row(string time, string request_type, string input, string output){
            output_file << "<tr><th>" << time << "</th>";
            output_file << "<th>" << request_type << "</th>";
            output_file << "<th>" << input << "</th>";
            output_file << "<th>" << output << "</th></tr>";
        }

        void close_file(){
            output_file << footer;
        }
};

int main(){
    HTMLLogger html_logger(".", "test.html", "DJANGO");
    html_logger.create_file();
    html_logger.add_row("01/02/2023, 12:00:05", "POST", "input:{wee:baguette}", "ERROR");
    html_logger.add_row("01/02/2023, 12:00:07", "POST", "input:{hi:hello}", "OK");
    html_logger.close_file();
}