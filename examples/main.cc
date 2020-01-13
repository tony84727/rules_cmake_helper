#include <iostream>
#include "config.h"

int main() {
    if (HAVE_A) {
        std::cout << "HAVE_A" << std::endl;
    }
    if (HAVE_B) {
        std::cout << "HAVE_B" << std::endl;
    }
}