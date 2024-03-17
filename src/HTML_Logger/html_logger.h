#ifndef HTMLLOGGER_H
#define HTMLLOGGER_H

#include <iostream>
#include <fstream>
#include <string>

using namespace std;

class HTMLLogger {
private:
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
