/**
 * Script.js
 */
$(document).ready(function () {
    // Calcul du path - ajout du path /hugo si le site est piloter par Victor
    function computePath(path) {
        // if (path.indexOf("http") > -1) {
        //     return path;
        // } else if (window.location.pathname.indexOf("hugo/") > -1) {
        //     return "/hugo" + path;
        // }
        return path
    }

    $('.bee-admin').each(function () {
        if ($BeeIsServer == 'true') {
            $(this).removeClass('bee-admin')
        } else if (window.location.pathname.indexOf("hugo/") > -1) {
            $(this).removeClass('bee-admin')
        }
    });

    $(".bee-open").on('click', function (event) {
        window.location = computePath($(this).data('path') + "/");
        event.preventDefault();
    });
    $(".bee-window-open").on('click', function (event) {
        $path = $(this).data('path')
        if ($path.indexOf("http") > -1 || $path.indexOf(".pdf") > -1) {
            window.open(computePath($path + "/"), "_blank")
        } else {
            window.location = computePath($path + "/");
        }
        event.preventDefault();
    });

    $(".bee-edit-open").on('click', function (event) {
        var $path = $(this).data('path');
        if (!$path) {
            $path = $BeeFilePath; // défini dans footer.html
        }
        if ($BeeIsServer == 'true') {
            $path = "http://localhost:8080/victor/document/" + $path;
        } else if (window.location.pathname.indexOf("hugo/")) {
            $path = "/victor/document/" + $path;
        }
        window.open($path, 'hugo-edit', 'left=' + (screen.availWidth - 1024 - 5) + ',top=5,width=1024,height=1014,scrolling=yes,scrollbars=yes,resizeable=yes');
        event.preventDefault();
    });

    $('.bee-select-download').on('click', function (event) {
        var $path = $(this).data('path');
        var $base = $(this).data('base');
        var link = document.createElement('a');
        link.href = $path;
        link.download = $base;
        link.click();
        // window.open($selected.paths, '_blank');
        event.preventDefault();
    });

    // CLIC IMAGE POPUP
    $('.bee-modal-image').on('click', function (event) {
        var $pdf = $(this).data('pdf');
        var $src = $(this).data('img');
        if (!$src) {
            $src = $(this).find('img').attr('src');
        } else {
            $src = computePath($src);
        }
        $('#bee-modal-image').find('img').attr('src', $src);
        if ($pdf) {
            $('#bee-modal-image').find('.approve').removeClass("bee-hidden");
            $pdf = computePath($pdf);
        }
        $('#bee-modal-image')
            .modal({
                closable: true,
                onHide: function () {
                    $('#bee-modal-image').find('.approve').addClass("bee-hidden");
                    return true;
                },
                onApprove: function () {
                    window.open($pdf, "_blank");
                }
            }).modal('show');
        event.preventDefault();
    });

    // PDF VIEWER
    $('.bee-modal-viewer').on('click', function (event) {
        var $path = $(this).data('path');
        var $base = $(this).data('base');
        var $type = $(this).data('type');
        var $height = screen.availHeight - 400;
        var $content = $('#bee-modal-viewer').find('.content');
        $content[0].innerHTML = '<object data="' + $path + '" type="' + $type + '" height="' + $height + '" width="100%" typemustmatch></object>';
        $('#bee-modal-viewer')
            .modal({
                closable: true,
                onHide: function () {
                    return true;
                },
                onApprove: function () {
                    var link = document.createElement('a');
                    link.href = $path;
                    link.download = $base;
                    link.click();
                }
            }).modal('show');
        event.preventDefault();
    });
    // TEXT VIEWER
    $('.bee-modal-text').on('click', function (event) {
        var $path = $(this).data('path');
        var $base = $(this).data('base');
        var $height = screen.availHeight - 400;
        $('#bee-text').height($height);
        $('#bee-text').load($path);
        $('#bee-modal-text')
            .modal({
                closable: true,
                onHide: function () {
                    return true;
                },
                onApprove: function () {
                    var link = document.createElement('a');
                    link.href = $path;
                    link.download = $base;
                    link.click();
                }
            }).modal('show');
        event.preventDefault();
    });

    // RETURN TO TOP
    $(window).scroll(function () {
        if ($(this).scrollTop() >= 50) {        // If page is scrolled more than 50px
            $('#return-to-top').fadeIn(200);    // Fade in the arrow
        } else {
            $('#return-to-top').fadeOut(200);   // Else fade out the arrow
        }
    });
    $('#return-to-top').click(function () {      // When arrow is clicked
        $('body,html').animate({
            scrollTop: 0                       // Scroll to top of body
        }, 500);
    });

    // MODAL PALYER
    $('.bee-modal-player').on('click', function (event) {
        // titre du morceau
        $('#bee-modal-player').find('p').html($(this).data('base'));
        // path du morceau
        $('.bee-player').data('path', $(this).data('path'));
        $('#bee-modal-player')
            .modal({
                closable: false,
                onDeny: function () {
                    $player.stop();
                    return true;
                },
                onVisible: function () {
                    return true;
                }
            }).modal('show');
        event.preventDefault();
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
            if ($path.indexOf("http") > -1) {
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
            $('.bee-player, .bee-radio').each(function () {
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

    $('.bee-radio').on('click', function (event) {
        $radio.click($(this));
        event.preventDefault();
    });
    var $radio = {
        audioElement: null,
        selector: null,
        path: null,
        isSourceLoaded: false,
        getPath: function (selector) {
            $path = selector.data('path')
            if ($path.indexOf("http") > -1) {
                return $path;
            } else if (window.location.pathname.indexOf("hugo/") > -1) {
                return "/hugo" + $path;
            }
            return $path;
        },
        init: function (selector) {
            this.selector = selector;
            this.path = this.getPath(selector);
            if (!this.isSourceLoaded) {
                this.loadSource();
            }
        },
        loadSource: function () {
            // console.log(this.path, 'source loading...');
            if (this.audioElement) {
                this.audioElement.src = this.path;
                this.audioElement.load();
            } else {
                this.audioElement = new Audio(this.path);
            }
            this.isSourceLoaded = true;
        },
        play: function () {
            if (this.isSourceLoaded) {
                this.audioElement.play();
            }
            this.uiPlay();
        },
        pause: function () {
            this.audioElement.pause();
            this.uiPause();
        },
        stop: function () {
            if (this.isSourceLoaded) {
                this.audioElement.pause();
            }
            this.uiInit();
        },
        uiWait: function () {
            // notched circle loading
            this.uiInit();
            this.selector.children('i').removeClass('file audio outline orange');
            this.selector.children('i').addClass('notched circle loading');
            this.selector.addClass('warning');
        },
        uiPause: function () {
            this.uiInit();
            this.selector.children('i').removeClass('play file audio outline orange');
            this.selector.children('i').addClass('play');
            this.selector.addClass('success');
        },
        uiPlay: function () {
            this.uiInit();
            this.selector.children('i').removeClass('play file audio outline orange');
            this.selector.children('i').addClass('pause');
            this.selector.addClass('error');
        },
        uiInit: function () {
            $('.bee-radio, .bee-player').each(function () {
                $(this).removeClass('success error warning');
                $(this).children('i').removeClass('pause play notched circle loading');
                $(this).children('i').addClass('file audio outline orange');
            })
        },
        isPlaying: function () {
            if (this.audioElement.paused) {
                return false
            } else {
                return true
            }
        },
        click: function (selector) {
            if (this.getPath(selector) == this.path) {
                // clic sur le player en cours
                if (!this.isSourceLoaded) {
                    return
                }
                // Pause ou Start du player
                if (this.isPlaying()) {
                    this.pause();
                } else {
                    this.play();
                }
            } else {
                // clic sur un nouveau player
                // arrêt du player en cours
                if (this.isSourceLoaded) {
                    // this.stop();
                }
                this.isSourceLoaded = false;
                // démarrage d'un nouveau player
                this.init(selector);
                this.play();
            } // end if CurrentPlayer
        },
    } // end $player

    /**
     * TOC: table des matières
     * <a class="ui label" href="lien"><i class="icone icon"></i>label</a>
     */
    if ($('#toc').length > 0) {
        var $main = $('#toc').closest('.main');
        var $html = "<p>";
        $main.children('h2').each(function (index) {
            $item = '<a class="ui label" href="#' + $(this).attr('id')
                + '"><i class="chevron circle right icon"></i>' + $(this).text() + '</a>';
            $html += $item;
        })
        $html += "</p>";
        $('#toc').html($html);
    }

    /**
     * Appel Masonry pour afficher les galeries
     * https://masonry.desandro.com/
     */
    // init Masonry
    if ($('#galerie').length > 0) {
        var $grid = $('#galerie').masonry({
            // options...
            itemSelector: '.bee-item-masonry'
        });
        // layout Masonry after each image loads
        $grid.imagesLoaded().progress(function () {
            $grid.masonry('layout');
        });
    }

    /**
     * SEMANTIC
     */
    $('.ui.accordion').accordion();
    $('#about').popup({ hoverable: true });
    ;
});