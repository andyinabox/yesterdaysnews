package clipstreamer

import "gitlab.com/andyinabox/yesterdaysnews/pkg/errorhandler"

func (s *Streamer) error(typ string, err error) {
	s.errs <- errorhandler.Err(typ, err)
}
