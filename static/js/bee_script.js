/**
 * Script.js
 */
$(document).ready(function () {

    // Sélection d'un dossier ou fichier
    $('.bee-click').on('click', function (event) {
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
            if ($(this).hasClass('bee-dblclick')) {
                $('.bee-button-edit').show();
            }
        }
        event.preventDefault();
    });
    // Ouverture d'un dossier ou fichier
    $('.bee-dblclick').on('dblclick', function (event) {
        window.location = $(this).data('action');
        event.preventDefault();
    });
    // Bouton Edit seulemnt sur markdown et image
    $(".bee-button-edit").on('click', function (event) {
        window.location = $('#bee-action').val();
        event.preventDefault();
    });

    $('.bee-submit').on('click', function (event) {
        var $form = $(this).closest('section').find('.form');
        // particularité pour l'ace-editor
        if (editor) {
            var $input = $form.find("input[name='document']");
            $input.val(editor.getValue());
        }
        $form.submit();
        event.preventDefault();
    });

    // ACTION RENAME
    $('.bee-modal-new').on('click', function (event) {
        var $form = $('#bee-modal-rename').find('form');
        $form.attr('action', $(this).data('action'));
        $('#bee-modal-rename').find('.header').html($(this).attr('title'));
        $('#bee-modal-rename').find("input[name='new_name']").val('');
        $('#bee-modal-rename')
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
        var $form = $('#bee-modal-rename').find('form');
        var $base = $('#bee-base').val();
        var $path = $('#bee-path').val();
        $form.attr('action', $(this).data('action') + $path);
        $('#bee-modal-rename').find('.header').html($(this).attr('title'));
        $('#bee-modal-rename').find("input[name='new_name']").val($base);
        $('#bee-modal-rename')
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
        var $path = $('#bee-path').val();
        $form.attr('action', $(this).data('action') + $path);
        $('#bee-modal-confirm').find('.header').html($(this).attr('title'));
        $('#bee-modal-confirm').find('.message>.header').html($path);
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
                url: '/api/folders',
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
        var $input = $(this).closest('form').find("input[name='image']");
        var $image = $(this).closest('form').find('img');
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
                console.log("onBeforeComplete", props);
                console.log("canvas-id", props.canvas.id);
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
                console.log("onComplete", props);
                return true;
            }
        });
        ImageEditor.open($url);
        event.preventDefault();
    });

    // ACE EDITOR
    var editor = null;
    $("#ace_editor").each(function (index) {
        var $form = $(this).closest('section').find('.form');
        var $input = $(this).closest('form').find("input[name='document']");
        // $(this)[0] pour récupérer le DOMElement
        editor = ace.edit($(this)[0]);
        // aff de l'éditeur
        editor.container.style.opacity = "";
        // def du language
        var ext = $(this).data("mode").replace('.', '');
        var mode = 'text';
        switch (ext) {
            case 'md':
                mode = "markdown";
                break;
            default:
                mode = ext
        }
        editor.session.setMode("ace/mode/" + mode);
        // editor.setAutoScrollEditorIntoView(true);
        // hauteur 
        editor.setOption("maxLines", 100);
        editor.setOption("theme", 'ace/theme/eclipse');
        editor.session.setUseWrapMode(true);
        editor.session.setTabSize(2);
        editor.setShowPrintMargin(false);
        editor.setReadOnly(false);
        $(this).css("fontSize", '13px');
        editor.commands.addCommand({
            name: 'Save',
            bindKey: { win: 'Ctrl-S', mac: 'Command-S' },
            exec: function (editor) {
                $input.val(editor.getValue());
                $form.submit();
            },
            readOnly: true // false if this command should not apply in readOnly mode
        });
        editor.session.on('change', function (delta) {
            // delta.start, delta.end, delta.lines, delta.action
            $(".bee-submit").removeClass('disabled');
            // $input.val(editor.getValue());
        });
    });

    // IHM SEMANTIC
    // $('.menu .item').tab();
    // $('.ui.checkbox').checkbox();
    // $('.ui.radio.checkbox').checkbox();
    // $('.ui.dropdown').dropdown();
    // $('select.dropdown').dropdown();
    $('.message .close')
        .on('click', function () {
            $(this)
                .closest('.message')
                .transition('fade')
                ;
        }
        );
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
     * OLD
     */

    // Calendar
    $('#standard_calendar')
        .calendar({
            ampm: false,
            text: {
                days: ['D', 'L', 'M', 'M', 'J', 'V', 'S'],
                months: ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Septembre', 'Octobre', 'Novembre', 'Decembre'],
                monthsShort: ['Jan', 'Fev', 'Mar', 'Avr', 'Mai', 'Juin', 'Juil', 'Aou', 'Sep', 'Oct', 'Nov', 'Dec'],
                today: 'Aujourd\'hui',
                now: 'Maintenant',
                am: 'AM',
                pm: 'PM'
            },
            // formatter: {
            //     date: function (date, settings) {
            //         if (!date) return '';
            //         var day = date.getDate();
            //         var month = date.getMonth() + 1;
            //         var year = date.getFullYear();
            //         return year + '-' + month + '-' + day;
            //     }
            // }
        });

    // CLIC IMAGE POPUP
    var $hugo_view = $('#hugo_view').val();
    var $hugo_refresh = $('#hugo_refresh').val();
    $('.hugo-modal-image').on('click', function (event) {
        var $src = $(this).data('src');
        $('#hugo-image').attr('src', $src)
        $('#hugo-modal-image')
            .modal({
                closable: true,
                onHide: function () {
                    isUsed = false;
                    return true;
                }
            }).modal('show');

        event.preventDefault();
    });

    // Coloriage syntaxique
    if ($("#codemirror-markdown").length != 0) {
        var myCodeMirror = CodeMirror.fromTextArea(
            document.getElementById('codemirror-markdown')
            , {
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

    // Collapse
    $('.crud-collapse').on('click', function (event) {
        var portlet = $(this).closest('div');
        if ($(this).hasClass('open')) {
            portlet.find('.icon').removeClass("open");
        } else {
            portlet.find('.icon').addClass("open");
        }
        portlet.find('.list').toggle();
        portlet.find('.message').toggle();
        // portlet.find('.content').toggle();
        event.preventDefault();
    });

    // CLIC URL
    $('.crud-jquery-url').on('click', function (event) {
        if (isUsed) {
            event.preventDefault();
            return
        }
        if (event.target.nodeName == "BUTTON") {
            // pour laisser la main à crud-jquery-button
            // Cas d'un button dans une card
            event.preventDefault();
            return
        }
        // Mémo du contexte dans un cookie
        if ($crud_view && $crud_view.length > 0) {
            Cookies.set($crud_view, this.id)
            $(this).addClass("crud-list-selected");
        }

        var $url = $(this).data('url');
        window.location = $url;
        event.preventDefault();
    });

    // CLIC BUTTON URL
    $('.crud-jquery-button').on('click', function (event) {
        var $target = $(this).data('target');
        if (!$target || $target == '') {
            window.location = $(this).data('url');
        } else {
            window.open($(this).data('url'), $target);
        }
        event.preventDefault();
    });

    // ACTION DEMANDE CONFIRMATION
    $('.crud-jquery-action').on('click', function (event) {
        var $url = $(this).data('url');
        if ($(this).data('confirm') == true) {
            $('#crud-action').html($(this).html());
            $('#crud-modal-confirm')
                .modal({
                    closable: false,
                    onDeny: function () {
                        return true;
                    },
                    onApprove: function () {
                        $('form').attr('action', $url);
                        $('form', document).submit();
                    }
                }).modal('show');
        } else {
            // Sans demande de confirmation
            $('form').attr('action', $url);
            $('form', document).submit()
        }
        event.preventDefault();
    });

    // CLIC IMAGE POPUP
    $('.crud-popup-image').on('click', function (event) {
        isUsed = true;
        // Mémo du contexte dans un cookie
        if ($crud_view && $crud_view.length > 0) {
            var $anchor = $(this).closest('.card');
            Cookies.set($crud_view, $anchor.attr('id'))
            $(this).closest('.cards').find('.crud-list-selected').removeClass('crud-list-selected');
            $anchor.addClass("crud-list-selected");
        }

        var $url = $(this).data('url');
        $('#crud-image').attr('src', $url)
        $('#crud-modal-image')
            .modal({
                closable: true,
                onHide: function () {
                    isUsed = false;
                    return true;
                }
            }).modal('show');
        event.preventDefault();
    });
    // CLIC IMAGE POPUP
    $('.crud-popup-chart').on('click', function (event) {
        isUsed = true;
        // Mémo du contexte dans un cookie
        if ($crud_view && $crud_view.length > 0) {
            var $anchor = $(this).closest('.card');
            Cookies.set($crud_view, $anchor.attr('id'))
            $(this).closest('.cards').find('.crud-list-selected').removeClass('crud-list-selected');
            $anchor.addClass("crud-list-selected");
        }

        var $html = $(this).html();
        var canvasParent = $('#crud-chart')
        canvasParent.html($html)
        $('#crud-modal-chart')
            .modal({
                closable: true,
                onHide: function () {
                    isUsed = false;
                    return true;
                },
                onVisible: function () {
                    drawChart(canvasParent.children("canvas"));
                }
            }).modal('show');
        event.preventDefault();
    });

    // SUPPRESSION D'UN ENREGISTREMENT
    $('.crud-jquery-delete').on('click', function (event) {
        $('#crud-modal-confirm')
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

    // APRES CHARGEMENT HTML ET JAVASCRIPT
    // CONTEXTE DE LA VUE
    var $crud_view = $('#crud_view').val()
    if ($crud_view && $crud_view.length > 0) {
        // Si recherche dans Cookie : aff du input et sélection
        var $search = $('#search').val();
        if ($search != "") {
            $('#crud-search-active').trigger('click');
        }
        // Positionnement sur la dernière ligne sélectionnée
        // voir ligne avec CrudIndexAnchor
        if (Cookies.get($crud_view)) {
            $anchor = $('#' + Cookies.get($crud_view))
            if ($anchor.length) {
                $('html, body').animate({
                    scrollTop: $anchor.offset().top - 100
                }, 1000)
                $anchor.addClass("crud-list-selected");
                // Collpase du folder
                if ($anchor.hasClass("message")) {
                    $anchorCollapse = $('.' + Cookies.get($crud_view))
                    $anchorCollapse.trigger("click");
                }
            }
        }
    }

    if ($hugo_view && $hugo_view.length > 0) {
        // Si recherche dans Cookie : aff du input et sélection
        var $search = $('#search').val();
        if ($search != "") {
            $('#crud-search-active').trigger('click');
        }
        // Positionnement sur la dernière ligne sélectionnée
        if (Cookies.get($hugo_view)) {
            var $cookie = Cookies.get($hugo_view);
            var $anchor = $('#' + $cookie)
            if ($anchor.length) {
                $('html, body').animate({
                    scrollTop: $anchor.offset().top - 100
                }, 1000)
                $anchor.addClass("crud-list-selected");
                // Collpase du folder
                var $folderid = $anchor.data("root");
                var $folder = $('#' + $folderid);
                $folder.trigger("click");
            }
        }
    }

    /**
     * Ouverture d'une fenêtre en popup
     * TODO voir si accepter par les browsers
     */
    $(document).on('click', '.hugo-window-open', function (event) {
        // Mémo du contexte dans un cookie
        if ($hugo_view && $hugo_view.length > 0) {
            var $anchor = $(this).closest('.message');
            Cookies.set($hugo_view, $anchor.attr('id'))
            $(this).closest('main').find('.crud-list-selected').removeClass('crud-list-selected');
            $anchor.addClass("crud-list-selected");
        }
        // Préparation window.open
        var height = $(this).data("height") ? $(this).data("height") : 'max';
        var width = $(this).data("width") ? $(this).data("width") : 'large';
        var posx = $(this).data("posx") ? $(this).data("posx") : 'left';
        var posy = $(this).data("posy") ? $(this).data("posy") : '3';
        var target = $(this).attr("target") ? $(this).attr("target") : 'hugo-win';
        if (window.opener == null) {
            window.open($(this).data('url')
                , target
                , computeWindow(posx, posy, width, height, false));
        } else {
            window.opener.open($(this).data('url')
                , target
                , computeWindow(posx, posy, width, height, false));
        }
        event.preventDefault();
    });


    /**
     * Fermeture de la fenêtre popup
     */
    $(document).on('click', '.crud-jquery-close', function (event) {
        if ($('#button_validate').length > 0
            && $('#button_validate').attr('disabled') != "disabled") {
            $('#crud-action').html("Abandonner les modifications ?");
            $('#crud-modal-confirm')
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