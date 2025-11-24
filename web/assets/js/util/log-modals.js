// Utilities for Access/Error log modals to keep index page lean and reduce merge conflicts
(function (window) {
  const NO_RECORD_TEXT = "No Record...";

  const accessLogColors = {
    timestamp: "#3c89e8",
    message: "#bcbcbc",
    blocked: "#e04141",
  };

  const errorLogColors = {
    levels: ["#3c89e8", "#008771", "#008771", "#f37b24", "#e04141", "#bcbcbc"],
    ids: ['#CADABF', '#5F6F65', '#FFDFD6', '#BC9F8B', '#C9DABF', '#9CA986', '#808D7C', '#E7E8D8', '#B5CFB7'],
  };

  function formatAccessLogs(logs = []) {
    let formattedLogs = '';

    logs.forEach((log, index) => {
      if (!log || log.length <= 3) {
        return;
      }

      const [date, time] = log.split(' ', 2);
      const messageStart = (date?.length ?? 0) + (time?.length ?? 0) + 2;
      const message = log.substr(messageStart);
      const messageColor = message.includes('-> blocked') ? accessLogColors.blocked : accessLogColors.message;

      if (index > 0) {
        formattedLogs += '<br/>';
      }

      formattedLogs += `<span style="color: ${accessLogColors.timestamp};">${date ?? ''} ${time ?? ''}</span> - `;
      formattedLogs += `<span style="color: ${messageColor};">${message}</span>`;
    });

    return formattedLogs;
  }

  function formatErrorLogs(logs = []) {
    let formattedLogs = '';
    const levels = ["DEBUG", "INFO", "NOTICE", "WARNING", "ERROR"];
    const levelsMap = { "[Debug]": levels[0], "[Info]": levels[1], "[Notice]": levels[2], "[Warning]": levels[3], "[Error]": levels[4] };

    let idColorIndex = 0;
    let lastLogId = '';

    logs.forEach((log, index) => {
      if (!log || log.length <= 3) {
        return;
      }

      const [date, time, levelTag, id] = log.split(' ', 4);

      if (!date || !time || !levelTag || !id) {
        if (index > 0) {
          formattedLogs += '<br/>';
        }
        formattedLogs += log;
        return;
      }

      const messageStart = date.length + time.length + levelTag.length + id.length + 4;
      const message = log.substr(messageStart);

      const level = levelsMap[levelTag];
      const levelIndex = levels.indexOf(level, levels) || 5;

      if (index > 0) {
        formattedLogs += '<br/>';
      }

      formattedLogs += `<span style="color: ${errorLogColors.levels[0]};">${date} ${time}</span> `;
      formattedLogs += `<span style="color: ${errorLogColors.levels[levelIndex]};">${level}</span> `;
      formattedLogs += ' - ';

      if (id.startsWith('[')) {
        if (lastLogId !== '' && lastLogId !== id) {
          idColorIndex++;
        }
        if (idColorIndex >= errorLogColors.ids.length) {
          idColorIndex = 0;
        }

        lastLogId = id;
        const idColor = errorLogColors.ids[idColorIndex];
        formattedLogs += `<span style="color: ${idColor};">${id}</span> <span style="color: ${errorLogColors.levels[5]};">${message}</span>`;
      } else {
        formattedLogs += `<span style="color: ${errorLogColors.levels[5]};">${id} ${message}</span>`;
      }
    });

    return formattedLogs;
  }

  function createAccessLogModal() {
    return {
      visible: false,
      logs: [],
      rows: 500,
      grep: '',
      loading: false,
      formattedLogs: '',
      show(logs) {
        this.visible = true;
        this.logs = Array.isArray(logs) ? logs : [];
        this.formattedLogs = this.logs.length > 0 ? formatAccessLogs(this.logs) : NO_RECORD_TEXT;
      },
      hide() {
        this.visible = false;
      },
    };
  }

  function createErrorLogModal() {
    return {
      visible: false,
      logs: [],
      rows: 500,
      grep: '',
      loading: false,
      formattedLogs: '',
      show(logs) {
        this.visible = true;
        this.logs = Array.isArray(logs) ? logs : [];
        this.formattedLogs = this.logs.length > 0 ? formatErrorLogs(this.logs) : NO_RECORD_TEXT;
      },
      hide() {
        this.visible = false;
      },
    };
  }

  async function openAccessLog(modal) {
    if (!modal) {
      return;
    }

    modal.loading = true;
    const msg = await HttpUtil.post(`/panel/api/server/access-log/${modal.rows}`, { grep: modal.grep });
    if (msg?.success) {
      modal.show(msg.obj);
    }
    await PromiseUtil.sleep(500);
    modal.loading = false;
  }

  async function openErrorLog(modal) {
    if (!modal) {
      return;
    }

    modal.loading = true;
    const msg = await HttpUtil.post(`/panel/api/server/error-log/${modal.rows}`, { grep: modal.grep });
    if (msg?.success) {
      modal.show(msg.obj);
    }
    await PromiseUtil.sleep(500);
    modal.loading = false;
  }

  window.LogModalService = {
    createAccessLogModal,
    createErrorLogModal,
    openAccessLog,
    openErrorLog,
  };
})(window);
