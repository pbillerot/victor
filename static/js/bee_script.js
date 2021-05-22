/**
 * Script.js
 */
jQuery(function () {

  // Chargement $file
  var $file = null
  $.ajax({
    url: '/victor/api/file'
  }).done(function (response) {
    $file = response.results;
    // Réglage de l'UI du formulaire des fichiers
    if ($file.Inline) {
      $('.bee-window-open').removeClass('bee-hidden');
    }
    if ($file.IsText || $file.IsImage || $file.IsDrawio || $file.IsMarkdown || $file.IsSystem) {
      $('#button_validate').removeClass('bee-hidden');
    }
  });

  // Chargement #bee-tree-folders
  $.ajax({
    url: '/victor/api/folders'
  }).done(function (response) {
    var $data = response.results
    // généralion du  tree en html
    var $html = "";
    // Ajout de la racine
    $html += '<div class="item"><i class="home icon"></i><div class="content"><a href="" class="header" data-path="/">...</a></div></div>';

    var $rang = 1;
    for (i = 0; i < $data.length; i++) {
      if ($data[i].rang > $rang) {
        $html += '<div class="list">'; // start list
      } else if ($data[i].rang < $rang) {
        for (ir = 0; ir < $rang - $data[i].rang; ir++) {
          $html += '</div>'; // end content
          $html += '</div>'; // end item
          $html += '</div>'; // end list
          $html += '</div>'; // end content
          $html += '</div>'; // end item
        }
      } else {
        if (i > 0) {
          $html += '</div>'; // end content
          $html += '</div>'; // end item
        }
      } // endif
      $html += '<div class="item">';
      $html += '<i class="folder outline icon"></i>';
      $html += '<div class="content">';
      if ($data[i].selected) {
        $html += '<a class="header selected" data-path="' + $data[i].path + '">'
          + $data[i].base + '</a>';
      } else {
        $html += '<a href="" class="header" data-path="' + $data[i].path + '">'
          + $data[i].base + '</a>';
      }
      $rang = $data[i].rang;
    }
    $html += '</div>'; // end content
    $html += '</div>'; // end item
    $html += '</div>'; // end list
    $('#bee-tree-folders').html($html);
    $('#bee-tree-folders .header').on('click', function (event) {
      $('#bee-tree-folders').find('.selected').removeClass('selected');
      $(this).addClass('selected');
      // Update champ input dest
      var $folder = $(this).data('path')
      $(this).closest('form').find('input[name="dest"]').val($folder)
      // le message dans la modal
      $(this).closest('form').find('.bee-input-dest').html($folder)

      event.preventDefault();
    });
  });

  // Ouverture d'un fichier pour modification
  $('.bee-file-edit').on('tap', function (event) {
    // Mode sélection unique
    var $action = $(this).data('action')
    // Ouverture de l'éditeur viewer dans une fenêtre séparée à droite
    var $height = 'max';
    var $width = 'large';
    var $posx = 'right';
    var $posy = '5';
    var $target = 'hugo-file';
    if (window.opener == null) {
      window.open($action, $target, computeWindow($posx, $posy, $width, $height, false));
    } else {
      window.opener.open($action, $target, computeWindow($posx, $posy, $width, $height, false));
    }
    event.preventDefault();
  });

  // Ouverture d'un dossier ou fichier ou sélection multiple
  $('.bee-tap').on('tap', function (event) {
    if ($bee_selector == false) {
      // AUDIO
      if ($(this).hasClass('bee-modal-player')) {
        // titre du morceau
        $('#bee-modal-player').find('p').html($(this).data('base'));
        // path du morceau
        $('.bee-player').data('path', '/content/' + $(this).data('path'));
        $('#bee-modal-player')
          .modal({
            closable: false,
            onDeny: function () {
              $player.stop();
              $player.path = null;
              return true;
            },
            onVisible: function () {
              return true;
            }
          }).modal('show');
        return;
      }
      // Mode sélection unique
      var $action = $(this).data('action')
      if ($action.indexOf('/folder') != -1) {
        window.location = $action;
      } else if ($(this).hasClass('bee-download')) {
        // window.open($action, "download");
        location.replace($action);
      } else {
        // Ouverture de l'éditeur viewer dans une fenêtre séparée à droite
        var $height = 'max';
        var $width = 'large';
        var $posx = 'right';
        var $posy = '5';
        var $target = 'hugo-file';
        if (window.opener == null) {
          window.open($action, $target, computeWindow($posx, $posy, $width, $height, false));
        } else {
          window.opener.open($action, $target, computeWindow($posx, $posy, $width, $height, false));
        }
      }
    } else {
      // Mode sélection multiple
      if ($(this).hasClass('bee-selected')) {
        // désélection d'un item
        $(this).removeClass('bee-selected');
      } else {
        // sélection ajout d'un item
        $(this).addClass("bee-selected");
      }
      var qselected = $('.bee-selected').length
      if (qselected > 0) {
        $('.bee-modal-rename').show();
        $('.bee-select-download').show();
      }
      if (qselected > 0) {
        $('.bee-press-visible').show();
        $('.bee-selector').html(qselected);
        $('.bee-press-hidden-mobile').each(function () {
          $(this).addClass('bee-hidden-mobile')
        });
      } else {
        $('.bee-press-visible').hide();
        $('.bee-selected').removeClass('bee-selected');
        $('.bee-selector').html('<i class="check icon"></i>');
        $('.bee-press-hidden-mobile').each(function () {
          $(this).removeClass('bee-hidden-mobile')
        })
      }
      if (qselected > 1) {
        $('.bee-modal-rename').hide();
        $('.bee-select-download').hide();
      }
    }
    event.preventDefault();
  });
  // SELECTION MULTIPLE
  var $bee_selector = false;
  $('.bee-selector').on('tap', function (event) {
    if ($bee_selector) {
      $(this).removeClass('teal');
      $bee_selector = false;
      $(this).html('<i class="check icon"></i>')
      $('.bee-selected').removeClass('bee-selected');
      $('.bee-press-visible').hide();
    } else {
      $(this).addClass('teal');
      $bee_selector = true;
      $(this).html($('.bee-selected').length);
    }
    event.preventDefault();
  });

  // Sélection d'un dossier ou fichier
  $('.bee-press').on('press', function (event) {
    if ($bee_selector) {
      event.preventDefault();
      return;
    }
    if ($(this).hasClass('bee-selected')) {
      // désélection
      $('.bee-press-visible').hide();
      $(this).removeClass('bee-selected');
      // Element à réafficher sur press et sur mobile
      $('.bee-press-hidden-mobile').each(function () {
        $(this).removeClass('bee-hidden-mobile')
      });
    } else {
      // sélection
      $(this).parent().find('.bee-selected').removeClass('bee-selected');
      $(this).addClass("bee-selected");
      $('.bee-press-visible').show();
      $('.bee-modal-rename').show();
      $('.bee-select-download').show();
      // Element à cacher sur press et sur mobile
      $('.bee-press-hidden-mobile').each(function () {
        $(this).addClass('bee-hidden-mobile')
      });
    }
    event.preventDefault();
  });

  $('.bee-submit').on('click', function (event) {
    var $form = $('form');
    $form.submit();
    event.preventDefault();
  });

  // ACTION NEW
  $('.bee-modal-new').on('click', function (event) {
    var $form = $('#bee-modal-new').find('form');
    $form.attr('action', $(this).data('action'));
    $('#bee-modal-new').find('.header').html($(this).attr('title'));
    $('#bee-modal-new').find("input[name='new_name']").val('');
    $('#bee-modal-new')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $form.submit();
        }
      }).modal('show');
    event.preventDefault();
  });
  // ACTION RENAME
  $('.bee-modal-rename').on('click', function (event) {
    var $form = $('#bee-modal-new').find('form');
    // valorisation de paths et bases
    $selected = getSelectedPathHtml();
    $form.attr('action', $(this).data('action') + $selected.paths);
    $('#bee-modal-new').find('.bee-modal-title').html($(this).attr('title'));
    $('#bee-modal-new').find("input[name='new_name']").val($selected.baseUnique);
    $('#bee-modal-new')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $form.submit();
        }
      }).modal('show');
    event.preventDefault();
  });
  // ACTION DOWNLOAD
  $('.bee-select-download').on('click', function (event) {
    // Recherche du fichier sélectionné qui sera unique
    $selected = getSelectedPathHtml();
    var link = document.createElement('a');
    link.href = '/content' + $selected.paths;
    link.download = $selected.baseUnique;
    link.click();
    // window.open($selected.paths, '_blank');
    event.preventDefault();
  });
  // ACTION CONFIRMATION
  $('.bee-modal-confirm').on('click', function (event) {
    var $form = $('#bee-modal-confirm').find('form');
    $('.bee-modal-title').html($(this).attr('title'));
    $form.attr('action', $(this).data('action'));
    if ($(this).data('message')) {
      $('#bee-modal-confirm').find('.message>.header').html($(this).data('message'));
    }
    $('#bee-modal-confirm')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $('form', document).submit();
        }
      }).modal('show');
    event.preventDefault();
  });
  // ACTION UPLOAD
  $('.bee-modal-upload').on('click', function (event) {
    var $form = $('#bee-modal-upload').find('form');
    $('#bee-modal-upload')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $form.submit();
        }
      }).modal('show');
    event.preventDefault();
  });
  $('#bee-upload-file').on('change', function () {
    var $files = $(this).get(0).files;
    var $html = "";
    for (var i = 0; i < $files.length; i++) {
      var $filename = $files[i].name.replace(/.*(\/|\\)/, '');
      $html += '<div class="ui teal label">' + $filename + '</div>'
    }
    $('#bee-files-selected').html($html);
  });
  // ACTION COPIER ou DEPLACER
  $('.bee-modal-move').on('click', function (event) {
    var $modal = $('#bee-modal-move')
    // titre
    $modal.find('.bee-modal-title').html($(this).attr('title'));
    var $form = $modal.find('form');
    // valorisation de paths et bases
    $selected = getSelectedPathHtml();
    // Le champ input des fichiers sources
    $form.find('input[name="paths"]').val($selected.paths)
    $form.find('.bee-input-paths').html($selected.bases)
    // Le champ input du répertoire destination par défaut
    var $folder = $('#bee-ctx').data('folder')
    $form.find('input[name="dest"]').val($folder)
    $form.find('.bee-input-dest').html($folder)
    // l'action à déclencher sur le serveur
    $form.attr('action', $(this).data('action'));
    $('#bee-modal-move')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $form.submit();
        }
      }).modal('show');
    event.preventDefault();
  });

  // ACTION SUPPRIMER
  $('.bee-modal-delete').on('click', function (event) {
    var $modal = $('#bee-modal-confirm')
    // titre
    $modal.find('.bee-modal-title').html($(this).attr('title'));
    var $form = $modal.find('form');
    // valorisation de paths et bases
    $selected = getSelectedPathHtml();
    // Le champ input des fichiers sources
    $form.find('input[name="paths"]').val($selected.paths)
    // l'action à déclencher sur le serveur
    $form.attr('action', $(this).data('action'));
    // le message dans la modal
    $form.find('.message>.header').html($selected.bases);
    $('#bee-modal-confirm')
      .modal({
        closable: false,
        onDeny: function () {
          return true;
        },
        onApprove: function () {
          $('form', document).submit();
        }
      }).modal('show');
    event.preventDefault();
  });

  // CLIC IMAGE EDITOR POPUP
  $('.bee-popup-image-editor').on('click', function (event) {
    var $url = $(this).data('src');
    var $form = $(this).closest('body').find('.form');
    var $input = $form.find("input[name='image']");
    var $image = $form.find('img');
    const config = {
      language: 'fr',
      tools: ['adjust', 'effects', 'filters', 'rotate', 'crop', 'resize', 'text'],
      colorScheme: 'dark',
      translations: {
        fr: {
          'toolbar.download': 'Valider'
        },
      }
    };
    var mime = $url.endsWith('.png') ? 'image/png' : 'image/jpeg';
    // https://github.com/scaleflex/filerobot-image-editor
    const ImageEditor = new FilerobotImageEditor(config, {
      onBeforeComplete: (props) => {
        // console.log("onBeforeComplete", props);
        // console.log("canvas-id", props.canvas.id);
        var canvas = document.getElementById(props.canvas.id);
        var dataurl = canvas.toDataURL(mime, 1);
        // update image du browser
        $image.attr('src', dataurl);
        // remplissage du imput pour le submit
        $input.val(dataurl);
        $(".bee-submit").removeClass('disabled');
        return false;
      },
      onComplete: (props) => {
        // console.log("onComplete", props);
        return true;
      }
    });
    ImageEditor.open($url);
    event.preventDefault();
  });


  // CODEMIRROR : coloration syntaxique et auto-complétion
  // https://codemirror.net/
  // déclaration de la foncion autocomplete pour les mots définis dans shortcodes.js
  CodeMirror.commands.autocomplete = function (cm) {
    CodeMirror.showHint(cm, CodeMirror.hint.shortcode, { list: $shortcodes });
  }
  CodeMirror.hint.shortcode = function (cm, options) {
    var list = options.list || [];
    var cursor = cm.getCursor();
    var currentLine = cm.getLine(cursor.line);
    var start = cursor.ch;
    var end = start;
    while (end < currentLine.length && /[\w$]+/.test(currentLine.charAt(end))) ++end;
    while (start && /[\w$]+/.test(currentLine.charAt(start - 1))) --start;
    var curWord = start != end && currentLine.slice(start, end);
    var regex = new RegExp('' + curWord, 'i');
    var result = {
      list: (!curWord ? list : list.filter(function (item) {
        if (typeof item == "object") {
          return item.displayText.match(regex);
        } else {
          return item.match(regex);
        }
      })).sort(),
      from: CodeMirror.Pos(cursor.line, start),
      to: CodeMirror.Pos(cursor.line, end)
    };
    return result;
  };
  // Activation de CODEMIRROR
  if ($("#bee-editor").length != 0) {
    var $mode = $("#bee-editor").data("mode");
    var myCodeMirror = CodeMirror.fromTextArea(
      document.getElementById('bee-editor'), {
      lineNumbers: true,
      lineWrapping: true,
      mode: $mode,
      readOnly: false,
      theme: 'eclipse',
      viewportMargin: 20,
    }
    );
    myCodeMirror.focus();
    myCodeMirror.setCursor({ line: $('#cursor_line').val(), ch: $('#cursor_ch').val() });
    myCodeMirror.setOption("extraKeys", {
      Tab: function (cm) {
        var spaces = Array(cm.getOption("indentUnit") + 1).join(" ");
        cm.replaceSelection(spaces);
      },
      "Ctrl-S": function (cm) {
        var cursor = cm.getCursor();
        $('#cursor_ch').val(cursor.ch);
        $('#cursor_line').val(cursor.line);
        $("#button_validate").trigger('click');
      },
      "Ctrl-Space": "autocomplete",
      "Ctrl-/": function (cm) {
        cm.toggleComment();
      }
    });
    myCodeMirror.on("change", function (cm) {
      var cursor = cm.getCursor();
      $('#cursor_ch').val(cursor.ch);
      $('#cursor_line').val(cursor.line);
      $(".bee-submit").removeClass('disabled');
    })
  }

  $('#bee-upload-file').simpleUpload({
    url: '/victor/upload',
    method: 'post',
    // maxFileNum: 4,
    // maxFileSize: 10 * 1024 * 1024, // Bytes
    dropZone: '#bee-dropzone',
    progress: '#bee-progress',
  }).on('upload:before', function (e, file, i) {
    $('#bee-progress').removeClass('bee-hidden');
  }).on('upload:after', function (e, file, i) {
    window.location.reload();
  });

  /**
   * Clic sur un player parmi les players
   * data-path="path du fichier" (http ou /path..)
   * data-loop="false" (true par défaut)
   */
  $('.bee-player').on('click', function (event) {
    $player.click($(this));
    event.preventDefault();
  });
  var $player = {
    selector: null,
    path: null,
    isContextLoaded: false,
    isSourceLoaded: false,
    isLoop: true,
    context: null,
    source: null,
    getPath: function (selector) {
      $path = selector.data('path')
      if (window.location.pathname.indexOf("http") > -1) {
        return $path;
      } else if (window.location.pathname.indexOf("hugo/") > -1) {
        return "/hugo" + $path;
      }
      return $path;
    },
    init: function (selector) {
      this.selector = selector;
      this.path = this.getPath(selector);
      if (!this.isContextLoaded) {
        window.AudioContext = window.AudioContext || window.webkitAudioContext;
        this.context = new AudioContext();
        this.isContextLoaded = true;
        if (selector.data('loop') == "false") {
          this.isLoop = false;
        }
      }
      // console.log(this.path, 'init ok');
    },
    loadSource: function () {
      // console.log(this.path, 'source loading...');
      this.isSourceLoaded = false;
      // Requête asynchrone sur le serveur
      var $request = new XMLHttpRequest();
      $request.open('GET', this.path, true);
      $request.responseType = 'arraybuffer';
      $request.onload = function () {
        // Nous sommes dans un événement -> pas d'utilisation de this
        $player.context.decodeAudioData($request.response, function (buffer) {
          $player.source = $player.context.createBufferSource();
          $player.source.buffer = buffer;
          $player.source.connect($player.context.destination);
          if ($player.isLoop) $player.source.loop = true;
          $player.isSourceLoaded = true;
          $player.source.start(0);
          $player.uiPlay();
        }, function (e) {
          console.log("Erreur lors du décodage des données audio ", e.err);
        });
      }
      this.uiWait();
      $request.send();
    },
    stop: function () {
      if (this.isSourceLoaded) {
        this.source.stop(0);
        this.source.disconnect(0);
        this.context.resume();
      }
      this.uiInit();
    },
    uiWait: function () {
      // notched circle loading
      this.uiInit();
      this.selector.children('i').removeClass('file audio outline orange');
      this.selector.addClass('warning');
      this.selector.children('i').addClass('notched circle loading');

    },
    uiPause: function () {
      this.uiInit();
      this.selector.children('i').removeClass('play file audio outline orange');
      this.selector.addClass('success');
      this.selector.children('i').addClass('play');
    },
    uiPlay: function () {
      this.uiInit();
      this.selector.children('i').removeClass('play file audio outline orange');
      this.selector.addClass('error');
      this.selector.children('i').addClass('pause');
    },
    uiInit: function () {
      $('.bee-player').each(function () {
        $(this).removeClass('success error warning');
        $(this).children('i').removeClass('pause play notched circle loading');
        $(this).children('i').addClass('file audio outline orange');
      })
    },
    click: function (selector) {
      if (this.getPath(selector) == this.path) {
        // clic sur le player en cours
        if (!this.isSourceLoaded) {
          return
        }
        // Pause ou Start du player
        if (this.context.state === 'running') {
          this.context.suspend();
          this.uiPause();
        } else if (this.context.state === 'suspended') {
          this.context.resume();
          this.uiPlay();
        }
      } else {
        // clic sur un nouveau player
        // arrêt du player en cours
        if (this.isSourceLoaded) {
          this.stop();
        }
        // démarrage d'un nouveau player
        this.init(selector);
        this.loadSource();
      } // end if CurrentPlayer
    },
  } // end $player

  // IHM SEMANTIC
  // $('.menu .item').tab();
  // $('.ui.checkbox').checkbox();
  // $('.ui.radio.checkbox').checkbox();
  $('.ui.dropdown.item').dropdown();
  // $('select.dropdown').dropdown();
  $('.message .close')
    .on('click', function () {
      $(this)
        .closest('.message')
        .transition('fade');
    });
  // $('.hide')
  //     .on('click', function () {
  //         $(this)
  //             .closest('.message')
  //             .transition('fade')
  //             ;
  //     }
  //     );

  // Toaster
  $('#toaster')
    .toast({
      class: $('#toaster').data('color'),
      position: 'bottom right',
      message: $('#toaster').val()
    });

  /**
   * Ouverture d'une fenêtre en popup
   */
  $(document).on('click', '.bee-window-open', function (event) {
    // Préparation window.open
    var height = $(this).data("height") ? $(this).data("height") : 'max';
    var width = $(this).data("width") ? $(this).data("width") : 'large';
    var posx = $(this).data("posx") ? $(this).data("posx") : 'left';
    var posy = $(this).data("posy") ? $(this).data("posy") : '3';
    var target = $(this).attr("target") ? $(this).attr("target") : 'hugo-win';
    if (window.opener == null) {
      window.open($(this).data('url'), target, computeWindow(posx, posy, width, height, false));
    } else {
      window.opener.open($(this).data('url'), target, computeWindow(posx, posy, width, height, false));
    }
    event.preventDefault();
  });

  /**
   * Fermeture de la fenêtre popup
   */
  $(document).on('click', '.bee-confirm-close', function (event) {
    if ($('#button_validate').length > 0 &&
      $('#button_validate').hasClass('disabled') == false) {
      $('#bee-modal-confirm')
        .modal({
          closable: false,
          onDeny: function () {
            return true;
          },
          onApprove: function () {
            window.close();
          }
        }).modal('show');
    } else {
      window.close();
    }
    event.preventDefault();
  });

  /**
   * retourne le HTML des chemins concaténés des fichiers et répertoires sélectionnés
   */
  function getSelectedPathHtml() {
    // valorisation de bee-path
    var $paths = ""; var $bases = ""; var $baseUnique = ""
    $('.bee-selected').each(function () {
      if ($paths.length > 0) {
        $paths += ",";
      }
      $paths += $(this).data('path');
      $bases += '<span class="ui teal label">' + $(this).data('base') + '</span>';
      $baseUnique = $(this).data('base');
    });
    return {
      paths: $paths, bases: $bases, baseUnique: $baseUnique
    }
  }

  // APPEL DRAWIO
  $('.bee-drawio').on('click', function (event) {
    var url = 'https://embed.diagrams.net/?embed=1&ui=atlas&spin=1&modified=unsavedChanges&proto=json';
    var source = $('#bee-drawio')[0];
    // var title = source.getAttribute('title')
    // url += '&title=' + title;
    if (source.drawIoWindow == null || source.drawIoWindow.closed) {
      // Implements protocol for loading and exporting with embedded XML
      var receive = function (evt) {
        if (evt.data.length > 0 && evt.source == source.drawIoWindow) {
          var msg = JSON.parse(evt.data);

          // Received if the editor is ready
          if (msg.event == 'init') {
            // Sends the data URI with embedded XML to editor
            source.drawIoWindow.postMessage(JSON.stringify(
              { action: 'load', xmlpng: source.getAttribute('src') }), '*');
          }
          // Received if the user clicks save
          else if (msg.event == 'save') {
            // Sends a request to export the diagram as XML with embedded PNG
            source.drawIoWindow.postMessage(JSON.stringify(
              { action: 'export', format: 'xmlpng', spinKey: 'saving' }), '*');
          }
          // Received if the export request was processed
          else if (msg.event == 'export') {
            // Updates the data URI of the image
            source.setAttribute('src', msg.data);
            $('input[name="image"]').val(msg.data);
            $(".bee-submit").removeClass('disabled');
          }

          // Received if the user clicks exit or after export
          if (msg.event == 'exit' || msg.event == 'export') {
            // Closes the editor
            window.removeEventListener('message', receive);
            source.drawIoWindow.close();
            source.drawIoWindow = null;
          }
        }
      };
      // Opens the editor
      window.addEventListener('message', receive);
      var $height = 'max';
      var $width = 'max';
      var $posx = '5';
      var $posy = '5';
      var $target = '_blank';
      source.drawIoWindow = window.open(url, $target, computeWindow($posx, $posy, $width, $height, false));
      // source.drawIoWindow = window.open(url);
    }
    else {
      // Shows existing editor window
      source.drawIoWindow.focus();
    }
  });
});

/**
 * Calcul du positionnement et de la taille de la fenêtre sur l'écran
 * @param {string} posx left center right ou px
 * @param {int} posy px
 * @param {string} pwidth max large xlarge ou px
 * @param {string} pheight max ou px
 * @param {srtien} full_screen yes no 0 1
 */
function computeWindow(posx, posy, pwidth, pheight, full_screen) {
  if (full_screen) {
    pheight = screen.availHeight - 70;
    pwidth = screen.availWidth - 6;
  }
  var height = pheight != null ? (/^max$/gi.test(pheight) ? screen.availHeight - 120 : pheight) : 830;
  var width = 900;
  if (pwidth != null) {
    width = pwidth;
    if (/^max$/gi.test(pwidth)) width = screen.availWidth - 6;
    if (/^large$|^l$/gi.test(pwidth)) width = 1024;
    if (/^xlarge$|^xl$/gi.test(pwidth)) width = 1248;
  } // end largeur
  var left = 3;
  if (posx != null) {
    left = posx;
    if (/^left$/gi.test(posx)) left = 3;
    if (/^right$/gi.test(posx)) left = screen.availWidth - width - 18;
    if (/^center$/gi.test(posx)) left = (screen.availWidth - width) / 2;
  } // end posx
  var top = 6
  if (posy != null) {
    height = screen.availHeight - posy - 10;
    top = posy;
  }

  return 'left=' + left + ',top=' + top + ',height=' + height + ',width=' + width + ',scrolling=yes,scrollbars=yes,resizeable=yes';
}

/**
 * Blocage du carriage return dans un champ input par exemple
 * @param {object event} event
 */
function blockCR(event) {
  if (event.keyCode == 13) {
    event.preventDefault();
  }
}