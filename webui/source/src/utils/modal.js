import Confirm from '@/components/Confirm';
import {Modal} from 'iview';

const showConfirm = (msg) => {
    Modal.confirm({
        width: '500px',
        scrollable: false,
        closable: true,
        render: (h) => {
          return h(Confirm, {
            props: {msg},
            on: {},
          });
        },
        // onOk: () => {},
      });
};

export {
  showConfirm,
};