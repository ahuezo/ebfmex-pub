$(document).ready(function(){       
	getcarrousel();
	pagina = 0;
	var queryVars = getVars();
	var hasVars = false;
	//if(queryVars.hasOwnProperty('word') || queryVars.hasOwnProperty('catMenu') || queryVars.hasOwnProperty('estadoMenu') || queryVars.hasOwnProperty('tipoMenu') )
	if(queryVars.hasOwnProperty('word')){
		hasVars = true;
		$("input[name=word]").val(queryVars.word);
	}
	if(queryVars.hasOwnProperty('catMenu')){
		hasVars = true;
		$("select[name=catMenu]").val(queryVars.catMenu);
	}
	if(queryVars.hasOwnProperty('estadoMenu')){
		hasVars = true;
		$("select[name=estadoMenu]").val(queryVars.estadoMenu);
	}
	if(queryVars.hasOwnProperty('tipoMenu')){
		hasVars = true;
		$("select[name=tipoMenu]").val(queryVars.tipoMenu);
	}
	$(".cargando").show();
	search();
	$("#buscarOferta").click(function(){                
		pagina = 0;    
		search();
		return false;
	});
	cargaOfertas = false;
        inSearch = false;
	$(document).scroll(function(){
		if(cargaOfertas && !inSearch){
			if(($(document).scrollTop() + $(window).height()) >= ($(document).height() - 10)){
				search();
			}
		}
	})

	$("#CloserLight").click(function() {
		$('#cuerpo').removeClass('noscroll');
		$('#lightback').addClass("hide"); 
		$('#lightfront').addClass("hide"); 
				  return false;
	});

	/*$(".lighter").click(function() {
		$('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
		$('#lightback').removeClass("hide"); 
		$('#lightfront').removeClass("hide"); 
				  return false;
	});*/
	});
	function getVars() {
	var delimiter = "?"; // using '#' here is great for AJAX apps.
	var separator = "&";
	var url = location.href;
	var get_exists = (url.indexOf(delimiter) > -1) ? true : false;
	if (get_exists) {
		var url_get = {};
		var params = url.substring(url.indexOf(delimiter)+1);
		var params_array = params.split(separator);
		for (var i=0; i < params_array.length; i++) {
			var param_name = params_array[i].substring(0,params_array[i].indexOf('='));
			var param_value = params_array[i].substring(params_array[i].indexOf('=')+1);
			url_get[param_name] = param_value;
		}
	  return url_get;
	}
	return false;
	}
	function lighterAjax(){
            $(".lighter").click(function() {
                    var id = $(this).parent().attr('id');
                    $.get('http://home.ebfmex-pub.appspot.com/wsdetalle',{id:id},function(data){
                        var empresa = JSON.parse(data);
                        console.log(empresa);
                    })
                    $('#cuerpo').addClass('noscroll');//importante ese impide que el fondo scrolle mientras la oferta si lo hace
                    $('#lightback').removeClass("hide");
                    $('#lightfront').removeClass("hide");
              return false;
            });
	}
	function getcarrousel() {
	$.get("/carr", "", function(response){
		$('#logo1').html(response);
	});
	}
	function search(){
            pagina ++;
            if(pagina == '1')
		  $(".ofertCont").html('')
            $(".ofertCont").append('<div class="cargando">Cargando...</div>');
            inSearch = true;
            $(".cargando").show();
            var keywords = ($("input[name=word]").val() == '¿Qué buscas?') ? '' : $("input[name=word]").val();
            var categoria = $("select[name=catMenu]").attr("value");
            var estado = $("select[name=estadoMenu]").attr("value");
            var tipo = $("select[name=tipoMenu]").attr("value");
            $.get("http://movil.ebfmxorg.appspot.com/search",{pagina:pagina, keywords:keywords, categoria:categoria, estado:estado , tipo:tipo, kind: 'Oferta'},function(data){
		console.log(pagina);		
		var ofertas = JSON.parse(data);
                $(".cargando").remove();
		if(ofertas.length >= 1){
		  cargaOfertas = true;
		  for(var i in ofertas){
			urlOferta = 'http://movil.ebfmex-pub.appspot.com/busqueda-de-ofertas.html';
			addOferta = '<div class="oferta bgWh" id="'+ofertas[i].IdOft+'">'
			addOferta += '<a href="#" class="lighter">'
			addOferta += '<span class="imgcont">'
			addOferta += '<img src="http://pruebas.ebfmex-pub.appspot.com'+ofertas[i].Logo+'" width="212" height="218" alt="'+ofertas[i].Oferta+'" title="'+ofertas[i].Oferta+'" />'
			addOferta += '</span>'
			addOferta += '<h3>'+ofertas[i].Oferta+'</h3>'
			addOferta += '</a>'
			addOferta += '<div class="col-30PR first" style="">'
			addOferta += '<a onClick="window.open(\'mailto:?subject=Conoce esta oferta&body=Conoce esta oferta de El buen fin ' + urlOferta +'\', this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer/sharer.php?u=http://pruebas.ebfmxorg.appspot.com/busqueda-de-ofertas.html">'
			addOferta += '<img src="../imgs/ofrtTemp/mtShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-40PR first" style="margin-top:5px;">'
			addOferta += '<a onClick="window.open(this.href, this.target, \'width=600,height=400\'); return false;" href="http://www.facebook.com/sharer.php?s=100&p[url]=' + urlOferta + '&p[images][0]=' + ofertas[i].Logo + '&p[title]= ' + ofertas[i].Oferta +'">'
			addOferta += '<img src="../imgs/ofrtTemp/fbShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '<div class="col-30PR first">'
			addOferta += '<a onClick="window.open(\'https://twitter.com/intent/tweet?text=Viendo \' + this.href, this.target, \'width=600,height=400\'); return false" href="' + urlOferta +'" class="btwitter" title="Compartelo en Twitter">'
			addOferta += '<img src="../imgs/ofrtTemp/twShare.jpg" alt="Compartir en Facebook" />'
			addOferta += '</a>'
			addOferta += '</div>'
			addOferta += '</div>';
			$(".ofertCont").append(addOferta);
			lighterAjax();
			//console.log(i + "_" + j + ': ' + ofertas[i][j] );
		  }    
		}else{
			cargaOfertas = false;
                        $(".ofertCont").append('No hay mas ofertas para esta busqueda');
		}
        $(".cargando").remove();
        inSearch = false;

	});
}

