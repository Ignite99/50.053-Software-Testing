#include <iostream>
#include <csignal>
#include <fstream>
#include <string>
#include <sstream>
#include <vector>
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
        cout << "@HTMLLogger: Header file not found!" << endl;
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
}

// create table headings and update column size
void HTMLLogger::create_table_headings(string style, const vector<string> &column_names){
    column_num = column_names.size();

    output_file << R"( <table> )" << endl;
    output_file << R"( <tr style=")" << style << R"(;">)" << endl;
    for(int i = 0; i < column_num; i++){
        output_file << R"(<th>)" << column_names[i] << R"(</th>)" << endl;
    }
    output_file << R"(</tr>)" << endl;
}

// add_row without style
void HTMLLogger::add_row(const vector<string> &row){
    int row_size = row.size();

    if(row_size != column_num){
        cout << "@HTMLLogger: Invalid number of columns!" << endl;
        return;
    }

    output_file << "<tr>" << endl;
    for(int i = 0; i < row_size; i++){
        output_file << "<th>" << row[i] << "</th>" << endl;
    }
    output_file << "</tr>" << endl;
}

// add_row with style
void HTMLLogger::add_row(string style, const vector<string> &row){
    int row_size = row.size();

    if(row_size != column_num){
        cout << "@HTMLLogger: Invalid number of columns!";
        return;
    }

    output_file << "<tr>" << endl;
    for(int i = 0; i < row_size; i++){
        output_file << R"(<th style=")" << style << R"(;">)" << row[i] << R"(</th>)" << endl;
    }
    output_file << "</tr>" << endl;
}

void HTMLLogger::close_file(){
    ifstream footer_file;

    footer_file.open(footerFilePath);
    if (!footer_file){
        cout << "@HTMLLogger: Footer file not found!" << endl;
        return;
    }

    string line;
    while (getline(footer_file, line)){
        output_file << line << '\n';
    }
    output_file.close();
}