<!-- list.html -->
{{$composter := .Composter}}
{{$search := .Search}}
{{$records := .Records}}
<!DOCTYPE html>
<html>
{{template "_head.html" .}}

<body>
  <!-- MENU HORIZONTAL FIXED -->
  <div class="ui top fixed menu" style="background-color: #f9fafb;">
    <div class="ui container">
      <a href="#" class="header item">
        <img class="logo" src="{{.Config.Icon}}">
        {{.Config.Title}}
      </a>
    </div>
  </div>
  <main>
    <div class="ui container">
      {{template "_flash.html" .}}
      {{ $color := "yellow"}}
      <div class="ui stackable cards">
      <!-- BOUCLE Records -->
      {{ range $irecord, $record := .Records }}
        {{ $color = "green"}}
        {{ if $record.Inline}}{{$color = "green"}}{{end}}
        {{ if $record.Planified}}{{$color = "blue"}}{{end}}
        {{ if $record.Expired}}{{$color = "warning"}}{{end}}
        {{ if eq $record.Draft "1"}}{{$color = "error"}}{{end}}

        {{ if eq $record.IsDir 1 }}
        <!-- REPERTOIRE -->
          <div class="ui horizontal card">
            <div class="image">
              <i class="folder icon"></i>
            </div>
            <div class="content">
              <div class="header">Cute Dog</div>
              <div class="meta">
                <span class="category">Animals</span>
              </div>
              <div class="description">
                <p></p>
              </div>
            </div>
          </div>
              <div class="item">
                <!-- REPERTOIRES -->
                <i id="{{$record.Key}}" data-root="{{$record.Root}}"
                  class="folder outline large brown icon link crud-collapse"></i>
                <div class="content">
                  <a {{if eq $record.Key $record.Root}}id="{{$record.ID}}" {{else}}id="{{$record.Key}}" {{end}}
                    data-root="{{$record.Root}}" class="header hugo-window-open"
                    data-url="/directory/{{$record.Key}}" data-posx="right" target="hugo-view">
                    {{$record.Base}}</a>
          
          {{ else }}
            <!-- else isDir -->
            <!-- start if extension -->
            {{ if or (eq $record.Ext ".png") (eq $record.Ext ".jpg") (eq $record.Ext ".svg")}}
            <!-- IMAGE -->
            <a id="{{$record.Key}}" data-root="{{$record.Root}}" class="ui compact message hugo-window-open"
              {{if eq $search ``}} style="display: none;" {{end}}
              data-url="/image/{{$record.Key}}" data-posx="right" target="hugo-view">
              <i class="image outline blue large icon"></i>
              <span>{{$record.Base}}</span>
            </a>
            {{ else if (eq $record.Ext ".pdf")}}
            <!-- PDF -->
            <a id="{{$record.Key}}" data-root="{{$record.Root}}" class="ui compact message hugo-window-open" {{if eq
              $search ``}} style="display: none;" {{end}} data-url="/pdf/{{$record.Key}}"
              data-posx="right" target="hugo-view">
              <i class="file pdf outline red large icon"></i>{{$record.Base}}</a>
            {{ else if (eq $record.Ext ".yaml")}}
            <!-- YAML -->
            <a id="{{$record.Key}}" data-root="{{$record.Root}}" class="ui compact brown message hugo-window-open"
              {{if eq $search ``}} style="display: none;" {{end}}
              data-url="/document/{{$record.Key}}" data-posx="right" target="hugo-view">
              <i class="file code outline brown large icon"></i>{{$record.Base}}</a>
            {{else}}
            <!-- DOCUMENT -->
            <div id="{{$record.Key}}" data-root="{{$record.Root}}" class="ui {{$color}} compact message " {{if eq
              $search ``}} style="display: none;" {{end}}>
              <i class="file alternate outline green large icon"></i>
              <a class="ui medium text hugo-window-open" data-url="/document/{{$record.Key}}"
                data-posx="right" target="hugo-view">
                <b class="searchable">{{$record.Base}}</b>
              </a>
              <span class="ui purple text">[{{$record.Title}}]</span>
              {{$arr := BeeSplit $record.Tags ","}}
              {{range $i, $item := $arr}}
              <div class="ui grey basic label"><b>#{{$item}}</b></div>
              {{end}}
              {{$arr := BeeSplit $record.Categories ","}}
              {{range $i, $item := $arr}}
              <div class="ui tag grey basic label"><b>{{$item}}</b></div>
              {{end}}
              {{if (eq $record.Draft "1") }}
              <span class="ui red text" title="brouillon"><i class="firstdraft icon"></i></span>
              {{end}}
              {{if or $record.DatePublish $record.DateExpiry }}
              <div class="ui yellow image label">
                <i class="clock outline icon"></i>
                {{if $record.DatePublish}}<span title="date publication">
                  {{if $record.Planified }}<span class="ui red text">{{$record.DatePublish}}</span>
                  {{else}}{{$record.DatePublish}}{{end}}
                </span>{{end}}
                {{if $record.DateExpiry}}
                <div class="detail" title="date expiration">
                  {{if $record.Expired }}<span class="ui red text">{{$record.DateExpiry}}</span>
                  {{else}}{{$record.DateExpiry}}{{end}}
                </div>
                {{end}}
              </div>
              {{end}}
              <!-- end $record.Date -->
            </div>
            {{end}}
            <!-- end test extension -->
          {{end}}
          <!-- end isDir-->
          {{end}}
        <!-- end range $records -->
      </div> <!-- end cards -->
      {{end }}
    </div> <!-- end container-->
  </main>
  <!-- Demande de confirmation de l'action' -->
  <div id="crud-modal-confirm" class="ui modal">
    <div class="header" id="crud-action">Texte à venir</div>
    <div class="content">
      <p>Veuillez le confirmer</p>
    </div>
    <div class="actions">
      <div class="ui cancel button">Annuler</div>
      <div class="ui approve button">Je confirme</div>
    </div>
  </div> <!-- end modal confirm -->

  <!-- Affichage d'une image en popup-->
  <div id="hugo-modal-image" class="ui modal">
    <div class="actions">
      <div class="ui cancel button">Fermer</div>
    </div>
    <div class="image content">
      <img id="hugo-image" class="ui large image" src="">
    </div>
  </div> <!-- end modal image -->

  <form method="POST" action="/list">
    {{ .xsrfdata }}
    <input type="hidden" id="crud-form-search" name="search">
    <input type="hidden" id="crud-form-searchstop" name="searchstop">
  </form>
  <!-- CONTEXTE -->
  <input type="hidden" id="hugo_view" value="hugo-view">

  {{template "_foot.html" .}}
</body>

</html>