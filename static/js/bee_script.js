/**
 * Script.js
 */
$(document).ready(function () {

    // Ouverture d'un dossier ou fichier
    $('.bee-tap').on('tap', function (event) {
        var $action = $(this).data('action')
        if ($action.indexOf('/folder') != -1) {
            window.location = $action;
        } else {
            // Préparation window.open
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
        event.preventDefault();
    });
    // Sélection d'un dossier ou fichier
    $('.bee-press').on('press', function (event) {
        var $action = $(this).data('action')
        var $base = $(this).data('base')
        var $path = $(this).data('path')
        if ($(this).hasClass('bee-selected')) {
            $('.bee-hidden').hide();
            $(this).removeClass('bee-selected');
            $('#bee-action').val('')
            $('#bee-base').val('')
        } else {
            $(this).parent().find('.bee-selected').removeClass('bee-selected');
            $(this).addClass("bee-selected");
            $('.bee-hidden').show();
            $('#bee-action').val($action)
            $('#bee-base').val($base)
            $('#bee-path').val($path)
            // boutons edit si markdown et image
            $('.bee-button-edit').hide();
            if ($(this).hasClass('bee-tap')) {
                $('.bee-button-edit').show();
            }
        }
        event.preventDefault();
    });
    // Bouton Edit seulement sur markdown et image
    $(".bee-button-edit").on('click', function (event) {
        window.location = $('#bee-action').val();
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
        var $base = $('#bee-base').val();
        var $path = $('#bee-path').val();
        $form.attr('action', $(this).data('action') + $path);
        $('#bee-modal-new').find('.header').html($(this).attr('title'));
        $('#bee-modal-new').find("input[name='new_name']").val($base);
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
    // ACTION CONFIRMATION
    $('.bee-modal-confirm').on('click', function (event) {
        var $form = $('#bee-modal-confirm').find('form');
        $('#bee-modal-confirm').find('.header').html($(this).attr('title'));
        if ($(this).data('message')) {
            // cas submit action
            $form.attr('action', $(this).data('action'));
            $('#bee-modal-confirm').find('.message>.header').html($(this).data('message'));
        } else {
            var $path = $('#bee-path').val();
            $form.attr('action', $(this).data('action') + $path);
            $('#bee-modal-confirm').find('.message>.header').html($path);
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
    $('#bee-input-file').on('change', function () {
        var $files = $(this).get(0).files;
        var $html = "";
        for (var i = 0; i < $files.length; i++) {
            var $filename = $files[i].name.replace(/.*(\/|\\)/, '');
            $html += '<div class="ui label">' + $filename + '</div>'
        }
        $('#bee-files-selected').html($html);
    });
    // ACTION DEPLACER
    $('.bee-modal-move').on('click', function (event) {
        var $form = $('#bee-modal-move').find('form');
        var $path = $('#bee-path').val();
        $form.attr('action', $(this).data('action') + $path);
        $('#bee-modal-move').find('.header').html($(this).attr('title'));
        $('#bee-modal-move').find('.message>.header').html($path);
        $('#bee-ajax-folders').dropdown({
            apiSettings: {
                url: '/victor/api/folders',
                cache: false,
                onResponse: (response) => {
                    // console.log(response);
                    return response
                }
            },
            saveRemoteData: false
        });
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

    // CLIC IMAGE EDITOR POPUP
    $('.bee-popup-image-editor').on('click', function (event) {
        var $url = $(this).data('src');
        var $form = $(this).closest('body').find('.form');
        var $input = $form.find("input[name='image']");
        var $image = $form.find('img');
        const config = {
            language: 'fr',
            tools: ['adjust', 'effects', 'filters', 'rotate', 'crop', 'resize', 'text'],
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

    // Coloriage syntaxique CODEMIRROR
    if ($("#bee-editor").length != 0) {
        var myCodeMirror = CodeMirror.fromTextArea(
            document.getElementById('bee-editor'), {
            lineNumbers: true,
            lineWrapping: true,
            mode: 'yaml-frontmatter',
            readOnly: false,
            theme: 'eclipse',
            viewportMargin: 20
        }
        );
        myCodeMirror.on("change", function (cm) {
            $(".bee-submit").removeClass('disabled');
        })
        // CTRL+S
        $(window).bind('keydown', function (event) {
            if (event.ctrlKey || event.metaKey) {
                switch (String.fromCharCode(event.which).toLowerCase()) {
                    case 's':
                        event.preventDefault();
                        $(".bee-submit").trigger('click');
                        break;
                }
            }
        });
        $("#bee-editor").focus();
    }

    // IHM SEMANTIC
    // $('.menu .item').tab();
    // $('.ui.checkbox').checkbox();
    // $('.ui.radio.checkbox').checkbox();
    $('.ui.dropdown').dropdown();
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


    // Coloriage syntaxique
    if ($("#codemirror-markdown").length != 0) {
        var myCodeMirror = CodeMirror.fromTextArea(
            document.getElementById('codemirror-markdown'), {
            lineNumbers: false,
            lineWrapping: true,
            mode: 'yaml-frontmatter',
            readOnly: false,
            theme: 'eclipse',
            viewportMargin: 20
        }
        );
        myCodeMirror.on("change", function (cm) {
            $('#button_validate').removeAttr('disabled');
        })
    }

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