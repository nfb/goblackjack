
async function newGame(){
  const round = await fetch("/new",{method:"post"});
  clearCards();
  renderGame(await round.json());
}

function clearCards(){
  var ucardcount = document.querySelector('#yourcards').children.length;
  for (var i=0; i<ucardcount; i++){
    document.querySelector('#yourcards').children[0].remove();
  }
  var dcardcount = document.querySelector('#dealercards').children.length;
  for (var i=0; i<dcardcount; i++){
    document.querySelector('#dealercards').children[0].remove();
  }
}


function renderGame(round){
  console.log(round);
  for (var i=0;i<round["PlayerHand"].length;i++){
    createCard(true,round["PlayerHand"][i]);
  }
  for (i=0;i<round["DealerHand"].length;i++){
    createCard(false,round["DealerHand"][i]);
  }

}

function createCard(player,inputcard){
  const card = document.createElement("li");
  var innertext = ""
  //card.id = id;
  card.arbfield = 0;
  card.className = "playercard";
  if (inputcard["PictureCard"]){
    innertext = inputcard["PictureType"] + " of " + inputcard["Suit"]
  } else {
    innertext = inputcard["NumericRank"] + " of " + inputcard["Suit"]
  } 

  const cardText = document.createTextNode(innertext);
  card.appendChild(cardText);
  if (player == true){
    document.querySelector("#yourcards").appendChild(card);
  } else {
    document.querySelector("#dealercards").appendChild(card);
  }

}
export { newGame, createCard }
