#include <cmath>
#include <iostream>
#include <iomanip>

using namespace std;

void foo(int i) {
	char* nums = new char[1024 * 1024];
	if (i < 100) {
		foo(i + 1);
	}
}


int main(){
	foo(0);

    int h, m, s;
    char ch, aorp;

    cin >> h >> ch >> m >> ch >> s >> aorp >> ch;
    h = (aorp == 'A') ? (h==12 ? 0 : h) : (h==12 ? 12 : h+12);

    cout << setw(2) << setfill('0') << h << ":"
         << setw(2) << setfill('0') << m << ":"
         << setw(2) << setfill('0') << s << endl;

    return 0;
}