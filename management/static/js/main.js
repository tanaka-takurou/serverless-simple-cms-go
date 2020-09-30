$(document).ready(function() {
  var token = window.localStorage.getItem("accessToken");
  if (!token) {
    location.href = App.url
  }
});
var OpenModal = function() {
  $('.large.modal').modal('show');
};
var CloseModal = function() {
  $('.large.modal').modal('hide');
};
var ShowCategory2 = function() {
  $("#category2").removeClass("force_hide");
  $("#show_cat2").addClass("force_hide");
};
var ShowCategory3 = function() {
  $("#category3").removeClass("force_hide");
  $("#show_cat3").addClass("force_hide");
};
var ClickMenu = function() {
  $('.ui.sidebar').sidebar('toggle');
};
var SetConst = function() {
  var title = $("#sitetitle").val();
  var image = $("#headimage").val();
  if (!title) {
    $("#set_const").removeClass('disabled');
    $("#warning").text("No title.").removeClass("hidden").addClass("visible");
    return false;
  }
  if (!image) {
    $("#set_const").removeClass('disabled');
    $("#warning").text("No image.").removeClass("hidden").addClass("visible");
    return false;
  }
  const data = {action: "set_const", token: "", title, image};
  actionHandle("#set_const", data, (res)=>{
    $("#set_const").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var SetSample = function() {
  const data = {action: "set_sample", token: ""};
  actionHandle("#set_sample", data, (res)=>{
    $("#set_sample").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var GetItemCategoryList = function() {
  const data = {action: "get_item_category_list", token: ""};
  actionHandle("#get_item_category_list", data, (res)=>{
    $("#info").removeClass("hidden").addClass("visible");
    App.itemList = res.itemList;
    res.itemList.forEach( function(item) {
      var d = JSON.parse(item.data)
      var optionTag = $("<option>" + d.title + "</option>", {
        "value": d.id
      });
      $("#selectItem").append(optionTag);
    });
    App.categoryList = res.categoryList;
    res.categoryList.forEach( function(item) {
      var optionTag = $("<option>", {
        "value": item.data
      });
      $("#category").append(optionTag);
    });
  }, onError);
}
var GetCategoryList = function() {
  const data = {action: "get_category_list", token: ""};
  actionHandle("#get_category_list", data, (res)=>{
    $("#info").removeClass("hidden").addClass("visible");
    App.categoryList = res.categoryList;
    res.categoryList.forEach( function(item) {
      var optionTag = $("<option>", {
        "value": item.data
      });
      $("#category").append(optionTag);
    });
  }, onError);
}
var AddItem = function() {
  const title = $('#formContent input[name="title"]').val();
  const description = $('#formContent input[name="description"]').val();
  const image = $('#formContent input[name="image"]').val();
  var categoryNames = [];
  const c1 = $('#formContent input[name="category1"]').val();
  const c2 = $('#formContent input[name="category2"]').val();
  const c3 = $('#formContent input[name="category3"]').val();
  if (!!c1 && c1.length > 0) {
    categoryNames.push(c1);
  }
  if (!!c2 && c2.length > 0) {
    categoryNames.push(c2);
  }
  if (!!c3 && c3.length > 0) {
    categoryNames.push(c3);
  }
  const categories = JSON.stringify(categoryNames);
  const data = {action: "add_item", token: "", title, description, image, categories};
  actionHandle("#add_item", data, (res)=>{
    $("#add_item").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var FixItem = function() {
  const id = $('#formContent input[name="id"]').val();
  const title = $('#formContent input[name="title"]').val();
  const description = $('#formContent input[name="description"]').val();
  const image = $('#formContent input[name="image"]').val();
  var categoryNames = [];
  const c1 = $('#formContent input[name="category1"]').val();
  const c2 = $('#formContent input[name="category2"]').val();
  const c3 = $('#formContent input[name="category3"]').val();
  if (!!c1 && c1.length > 0) {
    categoryNames.push(c1);
  }
  if (!!c2 && c2.length > 0) {
    categoryNames.push(c2);
  }
  if (!!c3 && c3.length > 0) {
    categoryNames.push(c3);
  }
  const categories = JSON.stringify(categoryNames);
  const old_categories = JSON.stringify(App.oldCategory)
  const data = {action: "fix_item", token: "", id, title, description, image, categories, old_categories};
  actionHandle("#fix_item", data, (res)=>{
    $("#fix_item").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var GetJs = function() {
  console.log("get_js")
}
var FixJs = function() {
  const data = {action: "fixjs", token: "", jsstring};
  actionHandle("#fix_css", data, (res)=>{
    $("#fix_css").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var GetCss = function() {
  console.log("get_css")
}
var FixCss = function() {
  const data = {action: "fix_css", token: "", cssstring};
  actionHandle("#fix_css", data, (res)=>{
    $("#fix_css").removeClass('disabled');
    $("#info").removeClass("hidden").addClass("visible");
  }, onError);
}
var actionHandle = function(element, data, callback, onerror) {
  $(element).addClass('disabled');
  data.token = window.localStorage.getItem("accessToken");
  if (!data.token) {
    $(element).removeClass('disabled');
    $("#warning").text("Not login.").removeClass("hidden").addClass("visible");
    return false;
  }
  request(data, callback, onerror);
}
var onError = function(e) {
  if (!!e.responseJSON) {
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  } else {
    console.log(e.message);
    $("#warning").text(e.message).removeClass("hidden").addClass("visible");
  }
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
function parseJson (data) {
  var res = {};
  for (i = 0; i < data.length; i++) {
    res[data[i].name] = data[i].value;
  }
  return res;
}
function toBase64 (file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = error => reject(error);
  });
}
function onConverted () {
  return function(v) {
    App.imgdata = v;
    $('#preview').attr('src', v);
  }
}
function UploadImage(elm) {
  if (!!App.imgdata) {
    $(elm).addClass("disabled");
    putImage();
  } else {
    CloseModal();
  }
}
function putImage() {
  var token = window.localStorage.getItem("accessToken");
  if (!token) {
    return false;
  }
  const file = $('#image').prop('files')[0];
  const data = {action: 'upload_img', filename: file.name, filedata: App.imgdata, token: token};
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    App.imgUrl = res.imgurl;
    if (App.imgUrl.length > 0) {
      $("#img_url").val(App.imgUrl);
    }
  })
  .fail(function(e) {
    console.log(e);
  })
  .always(function() {
    CloseModal();
  });
}
function ChangeImage () {
  const file = $('#image').prop('files')[0];
  toBase64(file).then(onConverted());
}
var SelectContent = function() {
  var i = $("#selectItem").prop("selectedIndex");
  if (i > 0) {
    var d = JSON.parse(App.itemList[i - 1].data);
    App.oldCategory = d.categoryids;
    $('#formContent input[name="id"]').val(App.itemList[i - 1].id);
    $('#formContent input[name="title"]').val(d.title);
    $('#formContent input[name="description"]').val(d.description);
    $('#formContent input[name="image"]').val(d.image);
    if (d.categoryids.length > 0) {
      $('#formContent input[name="category1"]').val(GetCategoryName(d.categoryids[0]));
    }
    if (d.categoryids.length > 1) {
      $('#formContent input[name="category2"]').val(GetCategoryName(d.categoryids[1]));
    }
    if (d.categoryids.length > 2) {
      $('#formContent input[name="category3"]').val(GetCategoryName(d.categoryids[2]));
    }
  }
};
var GetCategoryName = function(categoryId) {
  var category = App.categoryList.find(v => v.id == categoryId);
  if (!!category && category.data.length > 0) {
    return category.data;
  }
  return "";
};
var GetCategoryId = function(categoryName) {
  var category = App.categoryList.find(v => v.id == categoryId);
  if (!!category && category.data.length > 0) {
    return category.data;
  }
  return "";
};
var App = { imgdata: null, url: location.origin + {{ .ApiPath }}, imgUrl: '', itemList: null, categoryList: null, oldCategory: null };
