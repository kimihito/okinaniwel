module Main exposing (main)
import Browser
import Http
import Html exposing (Html, div, text, input, table, tr, td)
import Html.Events exposing(onInput)
import Json.Decode exposing (Decoder, field, string)
import Html.Attributes exposing (placeholder, value)

type alias Dog
    = { id : String
      , breed : String
      , place : String
      }
type Status = Failure | Loading | Success

-- MAIN

main : Program () Model Msg
main =
    Browser.element
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        }

-- MODEL
type alias Model =
  {
    status: Status
  , dogs : List Dog
  , input : String

  }

init : () -> (Model, Cmd Msg)
init _ =
    ({ status = Loading
    , dogs = []
    , input = ""
    }, getDogs)


-- UPDATE

type Msg
    = GotDogs (Result Http.Error (List Dog))
    | InputChange String

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
       GotDogs result ->
        case result of
            Ok dogs ->
                ({ model | dogs = dogs, status = Success }, Cmd.none)
            Err _ ->
                ({ model | status = Failure }, Cmd.none)

       InputChange text ->
        ( { model | input = text }, Cmd.none)

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- VIEW

view : Model -> Html Msg
view model =
    case model.status of
        Failure ->
            div []
                [ text "I could not load dogs." ]
        Loading ->
            text "loading..."
        Success ->
            let
                filteredDogs =
                    if model.input /= "" then
                        List.filter(\dog -> String.contains model.input dog.place) model.dogs
                    else
                        model.dogs
            in
            div []
                [ input [ placeholder "検索", value model.input, onInput InputChange] []
                , table []
                    (List.map viewDog filteredDogs)
                ]

viewDog : Dog -> Html Msg
viewDog dog =
    tr []
        [
             td []
                [ text dog.breed ]
        ,    td []
                [ text dog.place ]
        ,    td []
                [ text dog.id ]
        ]


-- HTTP
getDogs : Cmd Msg
getDogs =
    Http.get
        { url = "https://kimihito.github.io/okinaniwel/dogs.json"
        , expect = Http.expectJson GotDogs ( Json.Decode.list dogDecoder )
        }

dogDecoder : Decoder Dog
dogDecoder =
    Json.Decode.map3 Dog
        ( field "id" string)
        ( field "bread" string)
        ( field "place" string)