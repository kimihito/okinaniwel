module Main exposing (main)
import Browser
import Http
import Html exposing (Html, div, text, h2, table, tr, td)
import Json.Decode exposing (Decoder, field, string)

type alias Dog
    = { id : String
      , breed : String
      , place : String
      }

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
type Model
  = Failure
  | Loading
  | Success (List Dog)

init : () -> (Model, Cmd Msg)
init _ =
    (Loading, getDogs)


-- UPDATE

type Msg
    = GotDogs (Result Http.Error (List Dog))

update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
       GotDogs result ->
        case result of
            Ok dogs ->
                (Success dogs, Cmd.none)
            Err _ ->
                (Failure, Cmd.none)

-- SUBSCRIPTIONS

subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none

-- VIEW

view : Model -> Html Msg
view model =
    div []
        [ h2 [] [ text "dogs"]
        , viewDogs model
        ]

viewDogs : Model -> Html Msg
viewDogs model =
    case model of
        Failure ->
            div []
                [ text "I could not load dogs." ]

        Loading ->
            text "loading..."

        Success dogs ->
            table []
                (List.map viewDog dogs)

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