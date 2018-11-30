//----- socket.io -----//
var socket = io('', {transports: ['websocket']});//this domain
socket.on('added',function(data){
    alert("成功新增！")
    clear()
});
socket.on('error',function(data){
    alert("維修中...")
});

$('input[name="date"]').val(getToday());

var s = $('select')[0];
$('a').text(s.options[s.selectedIndex].value)
$('a').attr("href", "https://docs.google.com/spreadsheets/d/"+s.options[s.selectedIndex].value)

function addOutcome(){
    if($("input[name='item']").val()==""){alert("名稱不能為空");return;}
    if($("input[name='amount']").val()==""){alert("金額不能為空");return;}
    var data = $('form').serializeObject();
    if(data.payer=="")data.state="Y"
    else data.state="N"
    data.reimburse="N"
    data.outcome=data.amount
    data.income="0"
    //alert(JSON.stringify(data))
    socket.emit('add', data)
    /*
    socket.emit('add', {
      date  : "2018/11/27",
  		item  : "b-ball",
  	  payer : "Joseph",
  		state : "N",
  		reimburse : "Y",
  		income : "0",
  		outcome : "50",
    });
    */
}

function addIncome(){
    var data = $('form').serializeObject();
    data.state="Y"
    data.reimburse=""
    data.income=data.amount
    data.outcome="0"
    //alert(JSON.stringify(data))
    socket.emit('add', data)
}

function clear(){
    $('input[name="date"]').val(getToday());
    $('input[name="item"]').val("");
    $('input[name="payer"]').val("");
    $('input[name="amount"]').val("");
}

$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
};

function getToday(){
  var today = new Date();
  var dd = today.getDate();
  var mm = today.getMonth() + 1; //January is 0!

  var yyyy = today.getFullYear();
  if (dd < 10) {
    dd = '0' + dd;
  }
  if (mm < 10) {
    mm = '0' + mm;
  }
  var today = yyyy + '-' + mm + '-' + dd;
  return today
}
