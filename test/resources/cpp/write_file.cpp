#include <iostream>
#include <fstream>
#include <string>

using namespace std;

int main() {

  const string file_path = "/test.txt";

  ofstream fout;
  fout.open(file_path);
  fout.close();

  if (!fout) {
    cout << "write to file " << file_path << " failed" << endl;
  } else {
    cout << "ok" << endl;
  }


  return 0;
}
