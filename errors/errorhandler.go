package errors

import(
  "net/http"

  log "gopkg.in/inconshreveable/log15.v2"
)

type HttpError struct {
  Status      int
  Description string
}

func (h HttpError) Error() string {
  if h.Description == "" {
    return http.StatusText(h.Status)
  }
  return h.Description
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
  if err == nil {
    return
  }

  if err, ok := err.(HttpError); ok {
    log.Error("Error Status", "type" , err.Status)
    http.Error(w, err.Error(), err.Status)
    return
  }

  http.Error(w, "Sorry, an error occurred.", http.StatusInternalServerError)

  /*                                                                              
  if err, ok := err.(validationFailure); ok {                                   
          renderJSON(w, err, http.StatusBadRequest)                             
          return                                                                
  }                                                                             
                                                                                
  if isErrSqlNoRows(err) {                                                      
          http.NotFound(w, r)                                                   
          return                                                                
  }                                                                             
                                                                                
  logError(err)                                                                 
                                                                                
  http.Error(w, "Sorry, an error occurred", http.StatusInternalServerError)*/
}
