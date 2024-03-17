#include <iostream>
#include <csignal>
#include <fstream>
#include <string>
#include <sstream>
#include "html_logger.h"

using namespace std;

HTMLLogger::HTMLLogger(string output_file_path, string output_file_name, string project_type){
    this->output_file_path = output_file_path;
    this->output_file_name = output_file_name;
    this->project_type = project_type;
}

void HTMLLogger::create_file(){
    ifstream header_file;

    header_file.open(headerFilePath);
    if (!header_file){
        cout << "HTMLLogger Header file not found!" << endl;
        return;
    }

    output_file.open(output_file_path + output_file_name);

    // add opened header file into output file
    string line;
    while (std::getline(header_file, line)){
        output_file << line << '\n';
    }

    // append fuzz target heading
    if (project_type == "DJANGO"){
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
    cout << "here" << endl;
    output_file << "<tr><th>" << time << "</th>\n";
    output_file << "<th>" << request_type << "</th>\n";
    output_file << "<th>" << input << "</th>\n";
    output_file << "<th>" << output << "</th></tr>\n";
}

void HTMLLogger::close_file(){
    ifstream footer_file;

    footer_file.open(footerFilePath);
    if (!footer_file){
        cout << "HTMLLogger Footer file not found!" << endl;
        return;
    }

    string line;
    while (getline(footer_file, line)){
        output_file << line << '\n';
    }
    output_file.close();
}