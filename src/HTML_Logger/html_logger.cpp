#include <iostream>
#include <csignal>
#include <fstream>
#include <string>
#include <sstream>
#include <vector>
#include <algorithm>
#include "html_logger.h"

using namespace std;

HTMLLogger::HTMLLogger(string output_file_path, string output_file_name, string project_type){
    this->output_file_path = output_file_path;
    this->output_file_name = output_file_name;
    this->project_type = project_type;
}

void HTMLLogger::create_file(){
    output_file.open(output_file_name);
    if(!output_file){
        cout << "Creating file " << output_file_name << " ..." << endl;
        output_file.open(output_file_name);
    }
    output_file << header;
    if(project_type == "DJANGO"){
        output_file << "<p><b>Fuzz Target: </b>Django</p>";
    } else if (project_type == "COAP"){
        output_file << "<p><b>Fuzz Target: </b>CoAP</p>";
    } else if (project_type == "BLE"){
        output_file << "<p><b>Fuzz Target: </b>BLE</p>";
    }
    // TODO - display any other overall stats
    output_file << table_columns;
}

void HTMLLogger::add_row(string time, string request_type, string input, string output){
    output_file << "<tr><th>" << time << "</th>";
    output_file << "<th>" << request_type << "</th>";
    output_file << "<th>" << input << "</th>";
    output_file << "<th>" << output << "</th></tr>";
}

void HTMLLogger::close_file(){
    output_file << footer;
    output_file.close();
}

// int main(){
//     HTMLLogger html_logger(".", "test.html", "DJANGO");
//     html_logger.create_file();
//     html_logger.add_row("01/02/2023, 12:00:05", "POST", "input:{wee:baguette}", "ERROR");
//     html_logger.add_row("01/02/2023, 12:00:07", "POST", "input:{hi:hello}", "OK");
//     html_logger.close_file();
// }