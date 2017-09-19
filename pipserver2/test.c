#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "cfi.h"
#include "test.h"

/* Basic wrapper around cfi_rate that demonstrates the lifecycle or
   CFI creation, rating, response cracking and cleanup. We simply
   print the result to stdout; Returns true if successful, false
   otherwise. */

char Dummy[1024]="Dummy";

void *Init() {
     char *cfg_file="cfi.conf";
     
     return cfi_setup(cfg_file, stdout);
}


char *RateUrl(void *cfi_handle, const char *url) {
     char txt[1024];
     int  txt_sz;
     char cat[CFI_CATMAX];
     int rsvd;
     void *rsp;
     int okay;
     if((rsp = cfi_rate(cfi_handle, NULL,
                        url, strlen(url),
                        NULL, 0,
                        NULL, 0)) != NULL) {
          okay = 1;  
          /* Step 3: Collect rating data from the response    */

          /* If we looked up an URL, try display the found    */
          /* URL, also known as a FURL                        */
          txt_sz = sizeof(txt);
          if (cfi_get_furl(rsp, txt, &txt_sz) && txt_sz) {
               fprintf(stdout, "FURL: %.*s\n", txt_sz, txt);
          } else {
                    fprintf(stdout, "FURL: <none>\n");
          }

          /* Switch on the type of response and act accordingly */
          switch(cfi_get_stat(rsp)) {
          case CFI_RSP_NOCAT:
               fprintf(stdout, "rate: NO CATEGORY\n");
               break;
          case CFI_RSP_UNKNOWN:
               fprintf(stdout, "rate: UNKNOWN\n");
               break;
          case CFI_RSP_RATED:
               /* Iterate over and print categories with value */
               cat[0] = 0;
               while(cfi_get_cat(rsp, cat, sizeof(cat), &rsvd)) {
                    fprintf(stdout, "rate: RATED, '%s' '%d'\n", cat, rsvd);
               }
               break;
          case CFI_RSP_UNRATEABLE:
               fprintf(stdout, "rate: UNRATEABLE\n");
               okay = 0;
               break;
          case CFI_RSP_ERROR:
          default:
               txt_sz = sizeof(txt);
               /* Responses may have an error message describing */
               /* problems encountered while rating data         */
               if (cfi_get_emsg(rsp, txt, &txt_sz) && txt_sz) {
                    fprintf(stdout, "rate: ERROR: %s\n", txt);
               } else
                    fprintf(stdout, "rate: ERROR\n");
               okay = 0;
               break;
          }
          /* Step 4: Clean-up, we are done rating */
          cfi_free_response(&rsp); /* Free the response; Note: we are
                                      passing the address of the response. */
     }

     return  Dummy;
}
