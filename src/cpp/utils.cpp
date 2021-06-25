#include <hats.h>


/*******
* Structure:
* bit 39 - 32: Xaddress (in pixels)
* bit 31 - 24: Yaddress (in pixels)
* bit 23: Polarity (0 for OFF, 1 for ON)
* bit 22 - 0: Timestamp (in microseconds)
*
* Taken from Gochard's website, adapted from python.
* Made it less for-loopy
**/
stream preprocess(buffer in) {
    const int8_t NBYTES = 5;
    stream out;
    int x_address, y_address, p_address, t_address, multiple;
    x_address = ((1<<8)-1)<<32;
    y_address = ((1<<8)-1)<<24;
    p_address = (1<<23);
    t_address = (1<<23)-1;

    float unitconv = 1e-3;

    std::vector<int> group;
    for (int i = 0; i < in.size(); i+=NBYTES) {
        int bits = 0;
        int x, y, p, t;

        for (int j = 0; j < NBYTES; j++) {
            bits |= in[i+j];
            bits = bits<<8;
        }

        bits = bits>>8;
        x = (bits&x_address)>>32;
        y = (bits&y_address)>>24;
        p = (bits&p_address)>>23;
        t = (bits&t_address);

        if (y==240) {
            multiple++;
            continue;
        }

        t += multiple*(1<<13);

        event e = {};

        e.x = x;
        e.y = y;
        e.p = p;
        e.t = ((float) t)*unitconv;
        out.push_back(e);
    }
    return out;
}

void load_files() {
    int a;
}