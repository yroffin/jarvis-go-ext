#include <stdio.h>
#include <stdarg.h>
#include <dlfcn.h>

void hexDump (char *desc, void *addr, int len) {
    int i;
    unsigned char buff[17];
    unsigned char *pc = (unsigned char*)addr;

    // Output description if given.
    if (desc != NULL)
        printf ("%s:\n", desc);

    if (len == 0) {
        printf("  ZERO LENGTH\n");
        return;
    }
    if (len < 0) {
        printf("  NEGATIVE LENGTH: %i\n",len);
        return;
    }

    // Process every byte in the data.
    for (i = 0; i < len; i++) {
        // Multiple of 16 means new line (with line offset).

        if ((i % 16) == 0) {
            // Just don't print ASCII for the zeroth line.
            if (i != 0)
                fprintf(stderr, "  %s\n", buff);

            // Output the offset.
            fprintf(stderr, "  %04x ", i);
        }

        // Now the hex code for the specific character.
        fprintf(stderr, " %02x", pc[i]);

        // And store a printable ASCII character for later.
        if ((pc[i] < 0x20) || (pc[i] > 0x7e))
            buff[i % 16] = '.';
        else
            buff[i % 16] = pc[i];
        buff[(i % 16) + 1] = '\0';
    }

    // Pad out last line if not exactly 16 characters.
    while ((i % 16) != 0) {
        fprintf(stderr, "   ");
        i++;
    }

    // And print the final ASCII bit.
    fprintf(stderr, "  %s\n", buff);
}

struct spi_ioc_transfer {
        int           tx_buf;
        int           filler1;
        int           rx_buf;
        int           filler2;

        int           len;
        int           speed_hz;

        short           delay_usecs;
        char            bits_per_word;
        char            cs_change;
        char            tx_nbits;
        char            rx_nbits;
        short           pad;
};

int ioctl_debug(void * pointer)
{
  fprintf(stderr, "tx_buf: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->tx_buf, 8);
  fprintf(stderr, "rx_buf: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->rx_buf, 8);
  fprintf(stderr, "len: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->len, 4);
  fprintf(stderr, "speed_hz: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->speed_hz, 4);
  fprintf(stderr, "delay_usecs: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->delay_usecs, 2);
  fprintf(stderr, "bits_per_word: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->bits_per_word, 1);
  fprintf(stderr, "cs_change: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->cs_change, 1);
  fprintf(stderr, "tx_nbits: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->tx_nbits, 1);
  fprintf(stderr, "rx_nbits: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->rx_nbits, 1);
  fprintf(stderr, "pad: ");
  hexDump (NULL, (void *) &((struct spi_ioc_transfer *) pointer)->pad, 2);
}
