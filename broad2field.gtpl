<html>
    <head>
    <title></title>
    </head>
    <body>
      <script type="text/javascript">
       function popg(){
        
                  alert(valDateNumber("12345678"))
                }
     
      function  valDateNumber(numX){
                var numVal = 1;
                var i;
                var ch;
                if ((numX.length < 7) || (numX.length >7) ){
                        return false;
                }
             for(i=0;i<numX.length;i++){
                   switch(numX.charAt(i)){
                        case '0':break;
                        case  '1':
                        case  '2':
                        case  '3':
                        case  '4':
                        case  '5':
                        case  '6':
                        case  '7':
                        case  '8':
                        case  '9':
                              break;
                        default:return false;}
              }
     
            return true;
  
            }
        
        function check_info() {
             var companyName = document.getElementById('companyName').value;
             var caseNumber = document.getElementById('caseNumber').value;
             var validNumber = valDateNumber(caseNumber);
             if(validNumber  < 1){
              alert('Case Number requires numeric characters only and seven digits in length !!');
              return false;
             } else
             if (companyName == "" || casenumber == ""){
              alert('Empty Fields not allowed ' + validNumber);
              return false;
             } 
             return true;
        }

    </script>

        <form action="/notifiyField" method="post" onsubmit="return check_info();">
            Company Name:<input type="text" name="companyName" id="companyName">
            Case Number:<input type="text" name="caseNumber" id="caseNumber">
           <br/>
		 Severity level:  <select name="casePriority">
          <option value="P1">Serverity 1</option>
          <option value="P2">Serverity 2</option>
          <option value="P3">Serverity 3</option>
        </select>
         <br/>
            <input type="submit" value="Broadcast Alert">
        </form>
    </body>
</html>
