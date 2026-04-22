function findWordNumber(body, word, n) {
  const wordRegex = new RegExp(word, 'gi');
  const numRegex = new RegExp(n, 'g');
  
  let bestMatch = null;
  let bestDistance = Infinity;
  
  let wordMatch;
  while ((wordMatch = wordRegex.exec(body)) !== null) {
    // console.error('word match:', JSON.stringify(wordMatch[0]), 'at', wordMatch.index);
    const wordEnd = wordMatch.index + wordMatch[0].length;
    
    numRegex.lastIndex = wordEnd;
    let numMatch;
    if ((numMatch = numRegex.exec(body)) !== null) {
      // console.error('  num match:', numMatch[0], 'distance', numMatch.index - wordEnd);
      const distance = numMatch.index - wordEnd;
      if (distance < bestDistance) {
        bestDistance = distance;
        bestMatch = numMatch[0];
      }
    }      
  }
  
  return bestMatch;
}

function processEmail(body) {
  fn = findWordNumber(body, '[^а-яёА-ЯЁ]ФН[^а-яёА-ЯЁ]', '\\d{' + 16 + '}')
  fd = findWordNumber(body, '[^а-яёА-ЯЁ]ФД[^а-яёА-ЯЁ]', '\\d{' + 5 + '}\\d?')
  fpd = findWordNumber(body, '[^а-яёА-ЯЁ]ФПД?[^а-яёА-ЯЁ]', '\\d{' + 9 + '}\\d?')
  total = findWordNumber(body, '[^а-яёА-ЯЁ]ИТОГО?[^а-яёА-ЯЁ]', '\\d+[.]?\\d*')


  const dateRE = /(\d{2})\.(\d{2})\.(\d{4})/;
  const timeRE = /(\d{2}):(\d{2})/;

  const dateMatch = body.match(dateRE);
  const timeMatch = body.match(timeRE);

  console.log(fn, fd, fpd, total)

  if (!(fn && fd && fpd && total && dateMatch && timeMatch)) {
    return null
  }

  const [, day, month, year] = dateMatch;
  const [, hour, minute] = timeMatch;
  const timestamp = `${year}${month}${day}T${hour}${minute}`;

  return {
    "fn": fn,
    "fd": fd,
    "fpd": fpd,
    "timestamp": timestamp,
    "total": total,
  }
}

function findData(body) {
  const match = body.match(/t=\d+T\d+&s=[\d.]+&fn=\d+&i=\d+&fp=\d+&n=\d+/);
  if (match) {
      Logger.log("Found data: " + match[0]);
      return match[0]
  }
  data = processEmail(body)
  if (data) {
    return "t=" + data.timestamp + "&s=" + data.total + "&fn=" + data.fn + "&i=" + data.fd + "&fp=" + data.fpd + "&n=1"
  }
  return null
}

function processEmails() {
  Logger.log("Start")

  const secret = PropertiesService.getScriptProperties().getProperty('API_TOKEN');
  if (!secret) throw Error(`Secret is empty`)
  //before:01/01/2025 
  const threads = GmailApp.search("is:inbox -label:processed-by-receipt-api -label:filtered-by-receipt-api -label:failed-receipt-api -label:later-receipt-api (from:shop OR from:ofd OR from:taxcom OR from:beeline)", 0, 10);

  const processedLabel =   GmailApp.getUserLabelByName("processed-by-receipt-api") || GmailApp.createLabel("processed-by-receipt-api");
  const filteredLabel =   GmailApp.getUserLabelByName("filtered-by-receipt-api") || GmailApp.createLabel("filtered-by-receipt-api");
  const failedLabel =   GmailApp.getUserLabelByName("failed-receipt-api") || GmailApp.createLabel("failed-receipt-api");
  const laterLabel =   GmailApp.getUserLabelByName("later-receipt-api") || GmailApp.createLabel("later-receipt-api");


  Logger.log("Found threads")


  for (const thread of threads) {
    for (const message of thread.getMessages()) {
      Logger.log("Processing message from " + message.getFrom())
      const body = message.getBody(); // decoded HTML

      const data = findData(body)
        Logger.log(data)

      if (data) {

        response = UrlFetchApp.fetch("https://wasdetchan.online/receipts", {
          method: "POST",
          payload: {
            qrraw: data,          
          },
          headers: {
            Authorization: secret,
          },
          muteHttpExceptions: true,
        });

        const code = response.getResponseCode()
        const json = JSON.parse(response.getContentText());

        if (code != 200) {
          Logger.log("API fetch failed: " + json.status)
          if (json.code != 4) {
            thread.addLabel(failedLabel);
          } else {
            thread.addLabel(laterLabel);            
          }

        } else {
          Logger.log("API fetch success")
          thread.addLabel(processedLabel);
        }
      } else {
        // Logger.log(body)
        thread.addLabel(filteredLabel)
      } 
    }

  }
}

function deleteLabel() {
  var label = GmailApp.getUserLabelByName("processed-by-receipt-api");
  if (label) {
    label.deleteLabel();
    Logger.log("Deleted");
  } else {
    Logger.log("Label not found");
  }
  label = GmailApp.getUserLabelByName("filtered-by-receipt-api");
  if (label) {
    label.deleteLabel();
    Logger.log("Deleted");
  } else {
    Logger.log("Label not found");
  }
  label = GmailApp.getUserLabelByName("failed-receipt-api");
  if (label) {
    label.deleteLabel();
    Logger.log("Deleted");
  } else {
    Logger.log("Label not found");
  }
  label = GmailApp.getUserLabelByName("later-receipt-api");
  if (label) {
    label.deleteLabel();
    Logger.log("Deleted");
  } else {
    Logger.log("Label not found");
  }
}

