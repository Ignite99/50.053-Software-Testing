#ifndef HTMLLOGGER_H
#define HTMLLOGGER_H

#include <iostream>
#include <fstream>
#include <string>

using namespace std;

class HTMLLogger {
private:
    string headerFilePath = "./src/HTML_Logger/formats/header.html";
    string footerFilePath = "./src/HTML_Logger/formats/footer.html";
    string table_columns =
        R"(
            <table>
                <tr style="background-color:lightgrey;">
                    <th>Time</th>
                    <th>Request Type</th>
                    <th>Contents</th>
                    <th>Output</th>
                </tr>
        )";
    
    string project_type;
    string output_file_path;
    string output_file_name;
    ofstream output_file;

public:
    HTMLLogger(string output_file_path, string output_file_name, string project_type); // constructor
    void create_file();
    void add_row(string time, string request_type, string input, string output);
    void close_file();
};

#endif // HTMLLOGGER_H
