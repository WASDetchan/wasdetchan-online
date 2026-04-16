function findWordNumber(body, word, n) {
  const wordRegex = new RegExp(word, 'gi');
  const numRegex = new RegExp('\\d{' + n + '}', 'g');
  
  let bestMatch = null;
  let bestDistance = Infinity;
  
  let wordMatch;
  while ((wordMatch = wordRegex.exec(body)) !== null) {
    //console.error('word match:', JSON.stringify(wordMatch[0]), 'at', wordMatch.index);
    const wordEnd = wordMatch.index + wordMatch[0].length;
    
    numRegex.lastIndex = wordEnd;
    let numMatch;
    if ((numMatch = numRegex.exec(body)) !== null) {
      //console.error('  num match:', numMatch[0], 'distance', numMatch.index - wordEnd);
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
  fn = findWordNumber(body, '[^а-яёА-ЯЁ]ФН[^а-яёА-ЯЁ]', 16)
  fd = findWordNumber(body, '[^а-яёА-ЯЁ]ФД[^а-яёА-ЯЁ]', 5)
  fpd = findWordNumber(body, '[^а-яёА-ЯЁ]ФПД[^а-яёА-ЯЁ]', 10)

  console.log(fn, fd, fpd)

  if (fn && fd && fpd) {
    return {
      "fn": fn,
      "fd": fd,
      "fpd": fpd,
    }
  }
  else {
    return null
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
    Logger.log("Made data: " + data);
    return "t=20260101T0000&s=0&fn=" + data.fn + "&i=" + data.fd + "&fp=" + data.fpd + "n=1"
  }
  return null
}

function processEmails() {
  Logger.log("Start")
  console.log("Start1")

  const secret = PropertiesService.getScriptProperties().getProperty('API_TOKEN');
  if (!secret) throw Error(`Secret is empty`)
  
  const threads = GmailApp.search("is:inbox -label:processed-by-receipt-api -label:filtered- (from:shop OR from:ofd OR from:taxcom)", 0, 10);

  const processedLabel =   GmailApp.getUserLabelByName("processed-by-receipt-api") || GmailApp.createLabel("processed-by-receipt-api");
  const filteredLabel =   GmailApp.getUserLabelByName("filtered-by-receipt-api") || GmailApp.createLabel("processed-by-receipt-api");

  Logger.log("Found threads")


  for (const thread of threads) {
    for (const message of thread.getMessages()) {
      Logger.log("Processing message")
      const body = message.getBody(); // decoded HTML

      const data = findData(body)

      if (data) {

        response = UrlFetchApp.fetch("https://wasdetchan.online/receipts", {
          method: "POST",
          payload: {
            qrraw: data,          
          },
          headers: {
            Authorization: secret,
          },
        });
        if (response.status != 200) {
          Logger.log("API fetch failed: " + response)
        } else {
          Logger.log("API fetch success")
        }
        thread.addLabel(processedLabel);

      } else {
        thread.addLabel(filteredLabel)
      }
    }

  }
}
