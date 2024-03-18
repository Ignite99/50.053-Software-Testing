#ifndef HTMLLOGGER_H
#define HTMLLOGGER_H

#include <iostream>
#include <fstream>
#include <string>

using namespace std;

class HTMLLogger {
private:
    int column_num;
    string headerFilePath = "./src/HTML_Logger/formats/header.html";
    string footerFilePath = "./src/HTML_Logger/formats/footer.html";
    string project_type;
    string output_file_path;
    string output_file_name;
    ofstream output_file;

public:
    HTMLLogger(string output_file_path, string output_file_name, string project_type); // constructor
    void create_file();
    void create_table_headings(string style, const vector<string> &column_names);
    void add_row(const vector<string> &row);
    void add_row(string style, const vector<string> &row);
    
    void close_file();
};

#endif // HTMLLOGGER_H
