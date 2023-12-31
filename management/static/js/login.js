$(document).ready(function() {
  let params = new URLSearchParams(document.location.search.substring(1));
  let page = params.get("page");
  if (page == "profile") {
    getuser();
  } else if (page == "logout") {
    logout();
  }
});

var Login = function() {
  $("#submit").addClass('disabled');
  var name = $("#nickName").val();
  var pass = $('#password').val();
  if (!name | !pass) {
    return false;
  }
  const action = "login";
  const data = {action, name, pass};
  request(data, (res)=>{
    window.localStorage.setItem("accessToken", res.token);
    window.setTimeout(() => {location.href = "./?page=top";}, 1000);
  }, onError);
};

var ChangePass = function() {
  var token = window.localStorage.getItem("accessToken");
  if (!token) {
    return false;
  }
  if (!checkPass($("#newpassword").val())) {
    $("#newpassword").val('');
    $("#passwarning").removeClass("hidden").addClass("visible");
    return
  }
  $("#submit").addClass('disabled');
  var pass = $('#password').val();
  var newpass = $('#newpassword').val();
  if (!token | !pass | !newpass) {
    return false;
  }
  const action = "change_pass";
  const data = {action, token, pass, newpass};
  request(data, (res)=>{
    console.log(res);
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
};

var SignUp = function() {
  if (!checkPass($("#password").val())) {
    $("#password").val('');
    $("#passwarning").removeClass("hidden").addClass("visible");
    return
  }
  $("#submit").addClass('disabled');
  var mail = $("#email").val();
  var name = $("#nickName").val();
  var pass = $("#password").val();
  if (!mail | !name | !pass) {
    return false;
  }
  const action = "sign_up";
  const data = {action, mail, name, pass};
  request(data, (res)=>{
    console.log(res);
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
};

var Activate = function() {
  $("#submit").addClass('disabled');
  var name = $("#nickName").val();
  var code = $("#activationKey").val();
  if (!name | !code) {
      return false;
  }
  const action = "confirm_signup";
  const data = {action, name, code};
  request(data, (res)=>{
    console.log(res);
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
};

var GetUser = function() {
  var token = window.localStorage.getItem("accessToken");
  if (!token) {
    $("#settings").addClass("hidden");
    $("#warning").text("You are not logged in yet.").removeClass("hidden").addClass("visible");
    return false;
  }
  const action = "get_user";
  const data = {action, token};
  request(data, (res)=>{
    console.log(res);
    $("#name").text(res.name);
  }, onError);
};

var Logout = function() {
  var token = window.localStorage.getItem("accessToken");
  if (!token) {
    return false;
  }
  const action = "logout";
  const data = {action, token};
  request(data, (res)=>{
    console.log(res);
    window.localStorage.setItem("accessToken", "");
    console.log("return to top");
  }, onError);
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};

var checkPass = function(s) {
  const ra = new RegExp(/[0-9]+/);
  const rb = new RegExp(/[a-z]+/);
  const rc = new RegExp(/[A-Z]+/);
  const rd = new RegExp(/[#\\(\\)_\\-\\@\\%\\#\\&\\$\\^\\*]+/);
  return s.length > 7 && ra.test(s) && rb.test(s) && rc.test(s) && rd.test(s)
};

var onError = function(e) {
  console.log(e.responseJSON.message);
  $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  $("#submit").removeClass('disabled');
};
function OpenModal() {
  $('.large.modal').modal('show');
}
function CloseModal() {
  $('.large.modal').modal('hide');
}
var App = { imgdata: null, url: location.origin + {{ .ApiPath }}, imgUrl: '' };
