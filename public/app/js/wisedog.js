// Get Project list at upper right dropdown menu via AJAX call

$('#dropdown-target').click(function(){  
  // not expand menu yet
  if($(this).attr('aria-expanded') == 'false'){
    $.ajax({
      url : "/project/get/list?limit=10"
    }).done(function(data){
      var dataJSON = jQuery.parseJSON(data)
      var toBeInsert = "";
      for(i = 0; i < dataJSON.length; i++){
        toBeInsert += "<li class='prjs'><a href='/project/" + 
            dataJSON[i].Name + "'>" + dataJSON[i].Name + "</a></li>";
      }
      $('#spinner').before(toBeInsert);
      $('#spinner').addClass("gone");
    }).fail(function( ){
      var failed = "<li class='prjs'><i class='fa fa-warning'></i> Failed to load</li>";
      $('#spinner').before(failed);
      $('#spinner').addClass("gone");
    });
  }else{
    $('#spinner').removeClass("gone");
    $('.prjs').remove();
  }
});
