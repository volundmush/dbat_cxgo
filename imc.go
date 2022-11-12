package main

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"github.com/gotranspile/cxgo/runtime/stdio"
	"math"
	"unsafe"
)

const IMC_VERSION_STRING = "IMC2 Freedom CL-2.2 "
const IMC_VERSION = 2
const MAX_IMCHISTORY = 40
const MAX_IMCTELLHISTORY = 20
const IMC_DIR = "../imc/"
const IMC_BUFF_SIZE = 0x4000
const IMCCHAN_LOG = 1
const IMC_TELL = 1
const IMC_DENYTELL = 2
const IMC_BEEP = 4
const IMC_DENYBEEP = 8
const IMC_INVIS = 16
const IMC_PRIVACY = 32
const IMC_DENYFINGER = 64
const IMC_AFK = 128
const IMC_COLORFLAG = 256
const IMC_PERMOVERRIDE = 512
const LGST = 4096
const SMST = 1024

type imc_constates int

const (
	IMC_OFFLINE = imc_constates(iota)
	IMC_AUTH1
	IMC_AUTH2
	IMC_ONLINE
)

type imc_permissions int

const (
	IMCPERM_NOTSET = imc_permissions(iota)
	IMCPERM_NONE
	IMCPERM_MORT
	IMCPERM_IMM
	IMCPERM_ADMIN
	IMCPERM_IMP
)

type imc_channel struct {
	Next        *IMC_CHANNEL
	Prev        *IMC_CHANNEL
	Name        *byte
	Owner       *byte
	Operators   *byte
	Invited     *byte
	Excluded    *byte
	Local_name  *byte
	Regformat   *byte
	Emoteformat *byte
	Socformat   *byte
	History     [40]*byte
	Flags       int
	Level       int16
	Open        bool
	Refreshed   bool
}
type IMC_CHANNEL imc_channel
type imc_packet struct {
	First_data *IMC_PDATA
	Last_data  *IMC_PDATA
	From       [1024]byte
	To         [1024]byte
	Type       [1024]byte
	Route      [1024]byte
}
type IMC_PACKET imc_packet
type imc_packet_data struct {
	Next  *IMC_PDATA
	Prev  *IMC_PDATA
	Field [16384]byte
}
type IMC_PDATA imc_packet_data
type imc_siteinfo struct {
	Servername  *byte
	Rhost       *byte
	Network     *byte
	Serverpw    *byte
	Clientpw    *byte
	Localname   *byte
	Fullname    *byte
	Ihost       *byte
	Email       *byte
	Www         *byte
	Base        *byte
	Details     *byte
	Iport       int
	Minlevel    int
	Immlevel    int
	Adminlevel  int
	Implevel    int
	Rport       uint16
	Sha256      bool
	Sha256pass  bool
	Autoconnect bool
	Inbuf       [16384]byte
	Incomm      [16384]byte
	Outbuf      *byte
	Versionid   *byte
	Outsize     uint
	Outtop      int
	Desc        int
	State       uint16
}
type SITEINFO imc_siteinfo
type imc_remoteinfo struct {
	Next    *REMOTEINFO
	Prev    *REMOTEINFO
	Name    *byte
	Version *byte
	Network *byte
	Path    *byte
	Url     *byte
	Host    *byte
	Port    *byte
	Expired bool
}
type REMOTEINFO imc_remoteinfo
type imc_ban_data struct {
	Next *IMC_BAN
	Prev *IMC_BAN
	Name *byte
}
type IMC_BAN imc_ban_data
type imcchar_data struct {
	Imcfirst_ignore *IMC_IGNORE
	Imclast_ignore  *IMC_IGNORE
	Rreply          *byte
	Rreply_name     *byte
	Imc_listen      *byte
	Imc_denied      *byte
	Imc_tellhistory [20]*byte
	Email           *byte
	Homepage        *byte
	Aim             *byte
	Yahoo           *byte
	Msn             *byte
	Comment         *byte
	Imcflag         int
	Icq             int
	Imcperm         int
}
type IMC_CHARDATA imcchar_data
type imc_ignore struct {
	Next *IMC_IGNORE
	Prev *IMC_IGNORE
	Name *byte
}
type IMC_IGNORE imc_ignore
type imcucache_data struct {
	Next   *IMCUCACHE_DATA
	Prev   *IMCUCACHE_DATA
	Name   *byte
	Time   int64
	Gender int
}
type IMCUCACHE_DATA imcucache_data
type imc_color_table struct {
	Next   *IMC_COLOR
	Prev   *IMC_COLOR
	Name   *byte
	Mudtag *byte
	Imctag *byte
}
type IMC_COLOR imc_color_table
type imc_command_table struct {
	Next        *IMC_CMD_DATA
	Prev        *IMC_CMD_DATA
	First_alias *IMC_ALIAS
	Last_alias  *IMC_ALIAS
	Function    IMC_FUN
	Name        *byte
	Level       int
	Connected   bool
}
type IMC_CMD_DATA imc_command_table
type imc_help_table struct {
	Next  *IMC_HELP_DATA
	Prev  *IMC_HELP_DATA
	Name  *byte
	Text  *byte
	Level int
}
type IMC_HELP_DATA imc_help_table
type imc_cmd_alias struct {
	Next *IMC_ALIAS
	Prev *IMC_ALIAS
	Name *byte
}
type IMC_ALIAS imc_cmd_alias
type imc_packet_handler struct {
	Next *IMC_PHANDLER
	Prev *IMC_PHANDLER
	Func PACKET_FUN
	Name *byte
}
type IMC_PHANDLER imc_packet_handler
type who_template struct {
	Head      *byte
	Plrheader *byte
	Immheader *byte
	Plrline   *byte
	Immline   *byte
	Tail      *byte
	Master    *byte
}
type WHO_TEMPLATE who_template
type IMC_FUN func(ch *char_data, argument *byte)
type PACKET_FUN func(q *IMC_PACKET, packet *byte)

var imcwait int
var imcconnect_attempts int
var imc_sequencenumber uint
var imcpacketdebug bool = FALSE != 0
var default_packets_registered bool = FALSE != 0
var imcucache_clock int64
var imc_time int64
var imcperm_names [6]*byte = [6]*byte{libc.CString("Notset"), libc.CString("None"), libc.CString("Mort"), libc.CString("Imm"), libc.CString("Admin"), libc.CString("Imp")}
var this_imcmud *SITEINFO
var first_imc_channel *IMC_CHANNEL
var last_imc_channel *IMC_CHANNEL
var first_rinfo *REMOTEINFO
var last_rinfo *REMOTEINFO
var first_imc_ban *IMC_BAN
var last_imc_ban *IMC_BAN
var first_imcucache *IMCUCACHE_DATA
var last_imcucache *IMCUCACHE_DATA
var first_imc_color *IMC_COLOR
var last_imc_color *IMC_COLOR
var first_imc_command *IMC_CMD_DATA
var last_imc_command *IMC_CMD_DATA
var first_imc_help *IMC_HELP_DATA
var last_imc_help *IMC_HELP_DATA
var first_phandler *IMC_PHANDLER
var last_phandler *IMC_PHANDLER
var whot *WHO_TEMPLATE

func imcstrlcpy(dst *byte, src *byte, siz uint64) uint64 {
	var (
		d *byte  = dst
		s *byte  = src
		n uint64 = siz
	)
	if n != 0 && func() uint64 {
		p := &n
		*p--
		return *p
	}() != 0 {
		for {
			if (func() byte {
				p := func() *byte {
					p := &d
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()
				*func() *byte {
					p := &d
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}() = *func() *byte {
					p := &s
					x := *p
					*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
					return x
				}()
				return *p
			}()) == 0 {
				break
			}
			if func() uint64 {
				p := &n
				*p--
				return *p
			}() == 0 {
				break
			}
		}
	}
	if n == 0 {
		if siz != 0 {
			*d = '\x00'
		}
		for *func() *byte {
			p := &s
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}() != 0 {
		}
	}
	return uint64(int64(uintptr(unsafe.Pointer(s))-uintptr(unsafe.Pointer(src))) - 1)
}
func imcstrlcat(dst *byte, src *byte, siz uint64) uint64 {
	var (
		d    *byte  = dst
		s    *byte  = src
		n    uint64 = siz
		dlen uint64
	)
	for func() uint64 {
		p := &n
		x := *p
		*p--
		return x
	}() != 0 && *d != '\x00' {
		d = (*byte)(unsafe.Add(unsafe.Pointer(d), 1))
	}
	dlen = uint64(int64(uintptr(unsafe.Pointer(d)) - uintptr(unsafe.Pointer(dst))))
	n = siz - dlen
	if n == 0 {
		return dlen + uint64(C.strlen(s))
	}
	for *s != '\x00' {
		if n != 1 {
			*func() *byte {
				p := &d
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *s
			n--
		}
		s = (*byte)(unsafe.Add(unsafe.Pointer(s), 1))
	}
	*d = '\x00'
	return dlen + uint64(int64(uintptr(unsafe.Pointer(s))-uintptr(unsafe.Pointer(src))))
}
func imclog(format *byte, _rest ...interface{}) {
	var (
		buf [4096]byte
		ap  libc.ArgList
	)
	ap.Start(format, _rest)
	stdio.Vsnprintf(&buf[0], LGST, libc.GoString(format), ap)
	ap.End()
	basic_mud_log(libc.CString("[IMC] %s"), &buf[0])
}
func imcbug(format *byte, _rest ...interface{}) {
	var (
		buf [4096]byte
		ap  libc.ArgList
	)
	ap.Start(format, _rest)
	stdio.Vsnprintf(&buf[0], LGST, libc.GoString(format), ap)
	ap.End()
	basic_mud_log(libc.CString(" [IMC] ***BUG*** %s"), &buf[0])
}
func imcstrrep(src *byte, sch *byte, rep *byte) *byte {
	var (
		lensrc    int = int(C.strlen(src))
		lensch    int = int(C.strlen(sch))
		lenrep    int = int(C.strlen(rep))
		x         int
		y         int
		in_p      int
		newsrc    [4096]byte
		searching bool = FALSE != 0
	)
	newsrc[0] = '\x00'
	for func() int {
		x = 0
		return func() int {
			in_p = 0
			return in_p
		}()
	}(); x < lensrc; func() int {
		x++
		return func() int {
			p := &in_p
			x := *p
			*p++
			return x
		}()
	}() {
		if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == *sch {
			searching = TRUE != 0
			for y = 0; y < lensch; y++ {
				if *(*byte)(unsafe.Add(unsafe.Pointer(src), x+y)) != *(*byte)(unsafe.Add(unsafe.Pointer(sch), y)) {
					searching = FALSE != 0
				}
			}
			if searching {
				for y = 0; y < lenrep; func() int {
					y++
					return func() int {
						p := &in_p
						x := *p
						*p++
						return x
					}()
				}() {
					if in_p == (int(LGST - 1)) {
						newsrc[in_p] = '\x00'
						return &newsrc[0]
					}
					if *(*byte)(unsafe.Add(unsafe.Pointer(src), x-1)) == *sch {
						if *rep == '\x1b' {
							if y < lensch {
								if y == 0 {
									newsrc[in_p-1] = *(*byte)(unsafe.Add(unsafe.Pointer(sch), y))
								} else {
									newsrc[in_p] = *(*byte)(unsafe.Add(unsafe.Pointer(sch), y))
								}
							} else {
								y = lenrep
							}
						} else {
							if y == 0 {
								newsrc[in_p-1] = *(*byte)(unsafe.Add(unsafe.Pointer(rep), y))
							}
							newsrc[in_p] = *(*byte)(unsafe.Add(unsafe.Pointer(rep), y))
						}
					} else {
						newsrc[in_p] = *(*byte)(unsafe.Add(unsafe.Pointer(rep), y))
					}
				}
				x += lensch - 1
				in_p--
				searching = FALSE != 0
				continue
			}
		}
		if in_p == (int(LGST - 1)) {
			newsrc[in_p] = '\x00'
			return &newsrc[0]
		}
		newsrc[in_p] = *(*byte)(unsafe.Add(unsafe.Pointer(src), x))
	}
	newsrc[in_p] = '\x00'
	return &newsrc[0]
}
func imcone_argument(argument *byte, arg_first *byte) *byte {
	var (
		cEnd  int8
		count int
	)
	count = 0
	if arg_first != nil {
		*arg_first = '\x00'
	}
	if argument == nil || *argument == '\x00' {
		return nil
	}
	for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	cEnd = ' '
	if *argument == '\'' || *argument == '"' {
		cEnd = int8(*func() *byte {
			p := &argument
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	for *argument != '\x00' && func() int {
		p := &count
		*p++
		return *p
	}() <= math.MaxUint8 {
		if *argument == byte(cEnd) {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
			break
		}
		if arg_first != nil {
			*func() *byte {
				p := &arg_first
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &argument
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		} else {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
	}
	if arg_first != nil {
		*arg_first = '\x00'
	}
	for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	return argument
}
func imc_strip_colors(txt *byte) *byte {
	var (
		color *IMC_COLOR
		tbuf  [4096]byte
	)
	strlcpy(&tbuf[0], txt, LGST)
	for color = first_imc_color; color != nil; color = color.Next {
		strlcpy(&tbuf[0], imcstrrep(&tbuf[0], color.Imctag, libc.CString("")), LGST)
	}
	for color = first_imc_color; color != nil; color = color.Next {
		strlcpy(&tbuf[0], imcstrrep(&tbuf[0], color.Mudtag, libc.CString("")), LGST)
	}
	return &tbuf[0]
}
func color_itom(txt *byte, ch *char_data) *byte {
	var (
		color *IMC_COLOR
		tbuf  [4096]byte
	)
	if txt == nil || *txt == '\x00' {
		return libc.CString("")
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 8)) != 0 {
		strlcpy(&tbuf[0], txt, LGST)
		for color = first_imc_color; color != nil; color = color.Next {
			strlcpy(&tbuf[0], imcstrrep(&tbuf[0], color.Imctag, color.Mudtag), LGST)
		}
	} else {
		strlcpy(&tbuf[0], imc_strip_colors(txt), LGST)
	}
	return &tbuf[0]
}
func color_mtoi(txt *byte) *byte {
	var (
		color *IMC_COLOR
		tbuf  [4096]byte
	)
	if txt == nil || *txt == '\x00' {
		return libc.CString("")
	}
	strlcpy(&tbuf[0], txt, LGST)
	for color = first_imc_color; color != nil; color = color.Next {
		strlcpy(&tbuf[0], imcstrrep(&tbuf[0], color.Mudtag, color.Imctag), LGST)
	}
	return &tbuf[0]
}
func imc_to_char(txt *byte, ch *char_data) {
	var buf [8192]byte
	stdio.Snprintf(&buf[0], int(LGST*2), "%s\x1b[0m", color_itom(txt, ch))
	send_to_char(ch, libc.CString("%s"), &buf[0])
}
func imc_printf(ch *char_data, fmt *byte, _rest ...interface{}) {
	var (
		buf  [4096]byte
		args libc.ArgList
	)
	args.Start(fmt, _rest)
	stdio.Vsnprintf(&buf[0], LGST, libc.GoString(fmt), args)
	args.End()
	imc_to_char(&buf[0], ch)
}
func imc_to_pager(txt *byte, ch *char_data) {
	var buf [8192]byte
	stdio.Snprintf(&buf[0], int(LGST*2), "%s\x1b[0m", color_itom(txt, ch))
	imc_to_char(&buf[0], ch)
}
func imcpager_printf(ch *char_data, fmt *byte, _rest ...interface{}) {
	var (
		buf  [4096]byte
		args libc.ArgList
	)
	args.Start(fmt, _rest)
	stdio.Vsnprintf(&buf[0], LGST, libc.GoString(fmt), args)
	args.End()
	imc_to_pager(&buf[0], ch)
}
func imcstr_prefix(astr *byte, bstr *byte) bool {
	if astr == nil {
		imcbug(libc.CString("C.strncasecmp: null astr."))
		return TRUE != 0
	}
	if bstr == nil {
		imcbug(libc.CString("C.strncasecmp: null bstr."))
		return TRUE != 0
	}
	for ; *astr != 0; func() *byte {
		astr = (*byte)(unsafe.Add(unsafe.Pointer(astr), 1))
		return func() *byte {
			p := &bstr
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}()
	}() {
		if C.tolower(int(*astr)) != C.tolower(int(*bstr)) {
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func imccapitalize(str *byte) *byte {
	var (
		strcap [4096]byte
		i      int
	)
	for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(str), i)) != '\x00'; i++ {
		strcap[i] = byte(int8(C.tolower(int(*(*byte)(unsafe.Add(unsafe.Pointer(str), i))))))
	}
	strcap[i] = '\x00'
	strcap[0] = byte(int8(C.toupper(int(strcap[0]))))
	return &strcap[0]
}
func imc_hasname(list *byte, member *byte) bool {
	if list == nil || *list == '\x00' {
		return FALSE != 0
	}
	if C.strstr(list, member) == nil {
		return FALSE != 0
	}
	return TRUE != 0
}
func imc_addname(list **byte, member *byte) {
	var newlist [4096]byte
	if imc_hasname(*list, member) {
		return
	}
	if (*list) == nil || **(**byte)(unsafe.Add(unsafe.Pointer(list), unsafe.Sizeof((*byte)(nil))*0)) == '\x00' {
		strlcpy(&newlist[0], member, LGST)
	} else {
		stdio.Snprintf(&newlist[0], LGST, "%s %s", *list, member)
	}
	for {
		if (*list) != nil {
			libc.Free(unsafe.Pointer(*list))
			*list = nil
		}
		if true {
			break
		}
	}
	*list = C.strdup(&newlist[0])
}
func imc_removename(list **byte, member *byte) {
	var newlist [4096]byte
	if !imc_hasname(*list, member) {
		return
	}
	strlcpy(&newlist[0], imcstrrep(*list, member, libc.CString("")), LGST)
	for {
		if (*list) != nil {
			libc.Free(unsafe.Pointer(*list))
			*list = nil
		}
		if true {
			break
		}
	}
	*list = C.strdup(&newlist[0])
}
func imc_nameof(src *byte) *byte {
	var (
		name [1024]byte
		x    uint64
	)
	for x = 0; x < uint64(C.strlen(src)); x++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '@' {
			break
		}
		name[x] = *(*byte)(unsafe.Add(unsafe.Pointer(src), x))
	}
	name[x] = '\x00'
	return &name[0]
}
func imc_mudof(src *byte) *byte {
	var (
		mud    [1024]byte
		person *byte
	)
	if (func() *byte {
		person = C.strchr(src, '@')
		return person
	}()) == nil {
		strlcpy(&mud[0], src, SMST)
	} else {
		strlcpy(&mud[0], (*byte)(unsafe.Add(unsafe.Pointer(person), 1)), SMST)
	}
	return &mud[0]
}
func imc_channel_mudof(src *byte) *byte {
	var (
		mud [1024]byte
		x   uint64
	)
	for x = 0; x < uint64(C.strlen(src)); x++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == ':' {
			mud[x] = '\x00'
			break
		}
		mud[x] = *(*byte)(unsafe.Add(unsafe.Pointer(src), x))
	}
	return &mud[0]
}
func imc_channel_nameof(src *byte) *byte {
	var (
		name  [1024]byte
		x     uint64
		y     uint64 = 0
		colon bool   = FALSE != 0
	)
	for x = 0; x < uint64(C.strlen(src)); x++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == ':' {
			colon = TRUE != 0
			continue
		}
		if !colon {
			continue
		}
		name[func() uint64 {
			p := &y
			x := *p
			*p++
			return x
		}()] = *(*byte)(unsafe.Add(unsafe.Pointer(src), x))
	}
	name[x] = '\x00'
	return &name[0]
}
func imc_makename(person *byte, mud *byte) *byte {
	var name [1024]byte
	stdio.Snprintf(&name[0], SMST, "%s@%s", person, mud)
	return &name[0]
}
func escape_string(src *byte) *byte {
	var (
		newstr   [4096]byte
		x        uint64
		y        uint64 = 0
		quote    bool   = FALSE != 0
		endquote bool   = FALSE != 0
	)
	if C.strchr(src, ' ') != nil {
		quote = TRUE != 0
		endquote = TRUE != 0
	}
	for x = 0; x < uint64(C.strlen(src)); x++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '=' && quote {
			newstr[y] = '='
			newstr[func() uint64 {
				p := &y
				*p++
				return *p
			}()] = '"'
			quote = FALSE != 0
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '\n' {
			newstr[y] = '\\'
			newstr[func() uint64 {
				p := &y
				*p++
				return *p
			}()] = 'n'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '\r' {
			newstr[y] = '\\'
			newstr[func() uint64 {
				p := &y
				*p++
				return *p
			}()] = 'r'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '\\' {
			newstr[y] = '\\'
			newstr[func() uint64 {
				p := &y
				*p++
				return *p
			}()] = '\\'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(src), x)) == '"' {
			newstr[y] = '\\'
			newstr[func() uint64 {
				p := &y
				*p++
				return *p
			}()] = '"'
		} else {
			newstr[y] = *(*byte)(unsafe.Add(unsafe.Pointer(src), x))
		}
		y++
	}
	if endquote {
		newstr[func() uint64 {
			p := &y
			x := *p
			*p++
			return x
		}()] = '"'
	}
	newstr[y] = '\x00'
	return &newstr[0]
}
func imc_find_user(name *byte) *char_data {
	var (
		d   *descriptor_data
		vch *char_data = nil
	)
	for d = descriptor_list; d != nil; d = d.Next {
		if (func() *char_data {
			vch = func() *char_data {
				if d.Character != nil {
					return d.Character
				}
				return d.Original
			}()
			return vch
		}()) != nil && C.strcasecmp(GET_NAME(vch), name) == 0 && d.Connected == CON_PLAYING {
			return vch
		}
	}
	return nil
}
func imcgetname(from *byte) *byte {
	var (
		buf  [1024]byte
		mud  *byte
		name *byte
	)
	mud = imc_mudof(from)
	name = imc_nameof(from)
	if C.strcasecmp(mud, this_imcmud.Localname) == 0 {
		strlcpy(&buf[0], imc_nameof(name), SMST)
	} else {
		strlcpy(&buf[0], from, SMST)
	}
	return &buf[0]
}
func imc_isbanned(who *byte) bool {
	var mud *IMC_BAN
	for mud = first_imc_ban; mud != nil; mud = mud.Next {
		if C.strcasecmp(mud.Name, imc_mudof(who)) == 0 {
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func imc_isignoring(ch *char_data, ignore *byte) bool {
	var temp *IMC_IGNORE
	for temp = ch.Player_specials.Imcchardata.Imcfirst_ignore; temp != nil; temp = temp.Next {
		if C.strcasecmp(imc_nameof(temp.Name), libc.CString("*")) == 0 {
			if C.strcasecmp(imc_mudof(temp.Name), imc_mudof(ignore)) == 0 {
				return TRUE != 0
			}
		}
		if C.strcasecmp(imc_mudof(temp.Name), libc.CString("*")) == 0 {
			if C.strcasecmp(imc_nameof(temp.Name), imc_nameof(ignore)) == 0 {
				return TRUE != 0
			}
		}
		if !imcstr_prefix(ignore, temp.Name) {
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func imc_delete_info() {
	for {
		if this_imcmud.Servername != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Servername))
			this_imcmud.Servername = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Rhost != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Rhost))
			this_imcmud.Rhost = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Network != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Network))
			this_imcmud.Network = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Clientpw != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Clientpw))
			this_imcmud.Clientpw = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Serverpw != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Serverpw))
			this_imcmud.Serverpw = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Outbuf != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Outbuf))
			this_imcmud.Outbuf = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Localname != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Localname))
			this_imcmud.Localname = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Fullname != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Fullname))
			this_imcmud.Fullname = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Ihost != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Ihost))
			this_imcmud.Ihost = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Email != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Email))
			this_imcmud.Email = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Www != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Www))
			this_imcmud.Www = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Details != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Details))
			this_imcmud.Details = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Versionid != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Versionid))
			this_imcmud.Versionid = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud.Base != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Base))
			this_imcmud.Base = nil
		}
		if true {
			break
		}
	}
	for {
		if this_imcmud != nil {
			libc.Free(unsafe.Pointer(this_imcmud))
			this_imcmud = nil
		}
		if true {
			break
		}
	}
}
func imc_delete_reminfo(p *REMOTEINFO) {
	for {
		if p.Prev == nil {
			first_rinfo = p.Next
			if first_rinfo != nil {
				first_rinfo.Prev = nil
			}
		} else {
			p.Prev.Next = p.Next
		}
		if p.Next == nil {
			last_rinfo = p.Prev
			if last_rinfo != nil {
				last_rinfo.Next = nil
			}
		} else {
			p.Next.Prev = p.Prev
		}
		if true {
			break
		}
	}
	for {
		if p.Name != nil {
			libc.Free(unsafe.Pointer(p.Name))
			p.Name = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Version != nil {
			libc.Free(unsafe.Pointer(p.Version))
			p.Version = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Network != nil {
			libc.Free(unsafe.Pointer(p.Network))
			p.Network = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Path != nil {
			libc.Free(unsafe.Pointer(p.Path))
			p.Path = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Url != nil {
			libc.Free(unsafe.Pointer(p.Url))
			p.Url = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Port != nil {
			libc.Free(unsafe.Pointer(p.Port))
			p.Port = nil
		}
		if true {
			break
		}
	}
	for {
		if p.Host != nil {
			libc.Free(unsafe.Pointer(p.Host))
			p.Host = nil
		}
		if true {
			break
		}
	}
	for {
		if p != nil {
			libc.Free(unsafe.Pointer(p))
			p = nil
		}
		if true {
			break
		}
	}
}
func imc_new_reminfo(mud *byte, version *byte, netname *byte, url *byte, path *byte) {
	var (
		p        *REMOTEINFO
		mud_prev *REMOTEINFO
	)
	for {
		if (func() *REMOTEINFO {
			p = new(REMOTEINFO)
			return p
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	p.Name = C.strdup(mud)
	if url == nil || *url == '\x00' {
		p.Url = C.strdup(libc.CString("Unknown"))
	} else {
		p.Url = C.strdup(url)
	}
	if version == nil || *version == '\x00' {
		p.Version = C.strdup(libc.CString("Unknown"))
	} else {
		p.Version = C.strdup(version)
	}
	if netname == nil || *netname == '\x00' {
		p.Network = C.strdup(this_imcmud.Network)
	} else {
		p.Network = C.strdup(netname)
	}
	if path == nil || *path == '\x00' {
		p.Path = C.strdup(libc.CString("UNKNOWN"))
	} else {
		p.Path = C.strdup(path)
	}
	p.Expired = FALSE != 0
	for mud_prev = first_rinfo; mud_prev != nil; mud_prev = mud_prev.Next {
		if C.strcasecmp(mud_prev.Name, mud) >= 0 {
			break
		}
	}
	if mud_prev == nil {
		for {
			if first_rinfo == nil {
				first_rinfo = p
				last_rinfo = p
			} else {
				last_rinfo.Next = p
			}
			p.Next = nil
			if first_rinfo == p {
				p.Prev = nil
			} else {
				p.Prev = last_rinfo
			}
			last_rinfo = p
			if true {
				break
			}
		}
	} else {
		for {
			p.Prev = mud_prev.Prev
			if mud_prev.Prev == nil {
				first_rinfo = p
			} else {
				mud_prev.Prev.Next = p
			}
			mud_prev.Prev = p
			p.Next = mud_prev
			if true {
				break
			}
		}
	}
}
func imc_find_reminfo(name *byte) *REMOTEINFO {
	var p *REMOTEINFO
	for p = first_rinfo; p != nil; p = p.Next {
		if C.strcasecmp(name, p.Name) == 0 {
			return p
		}
	}
	return nil
}
func check_mud(ch *char_data, mud *byte) bool {
	var r *REMOTEINFO = imc_find_reminfo(mud)
	if r == nil {
		imc_printf(ch, libc.CString("~W%s ~cis not a valid mud name.\r\n"), mud)
		return FALSE != 0
	}
	if r.Expired {
		imc_printf(ch, libc.CString("~W%s ~cis not connected right now.\r\n"), r.Name)
		return FALSE != 0
	}
	return TRUE != 0
}
func check_mudof(ch *char_data, mud *byte) bool {
	return check_mud(ch, imc_mudof(mud))
}
func get_imcpermvalue(flag *byte) int {
	var x uint64
	for x = 0; x < uint64(unsafe.Sizeof([6]*byte{})/unsafe.Sizeof((*byte)(nil))); x++ {
		if C.strcasecmp(flag, imcperm_names[x]) == 0 {
			return int(x)
		}
	}
	return -1
}
func imccheck_permissions(ch *char_data, checkvalue int, targetvalue int, enforceequal bool) bool {
	if checkvalue < 0 || checkvalue > IMCPERM_IMP {
		imc_to_char(libc.CString("Invalid permission setting.\r\n"), ch)
		return FALSE != 0
	}
	if checkvalue > ch.Player_specials.Imcchardata.Imcperm {
		imc_to_char(libc.CString("You cannot set permissions higher than your own.\r\n"), ch)
		return FALSE != 0
	}
	if checkvalue == ch.Player_specials.Imcchardata.Imcperm && ch.Player_specials.Imcchardata.Imcperm != IMCPERM_IMP && enforceequal {
		imc_to_char(libc.CString("You cannot set permissions equal to your own. Someone higher up must do this.\r\n"), ch)
		return FALSE != 0
	}
	if ch.Player_specials.Imcchardata.Imcperm < targetvalue {
		imc_to_char(libc.CString("You cannot alter the permissions of someone or something above your own.\r\n"), ch)
		return FALSE != 0
	}
	return TRUE != 0
}
func imc_newban() *IMC_BAN {
	var ban *IMC_BAN
	for {
		if (func() *IMC_BAN {
			ban = new(IMC_BAN)
			return ban
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	ban.Name = nil
	for {
		if first_imc_ban == nil {
			first_imc_ban = ban
			last_imc_ban = ban
		} else {
			last_imc_ban.Next = ban
		}
		ban.Next = nil
		if first_imc_ban == ban {
			ban.Prev = nil
		} else {
			ban.Prev = last_imc_ban
		}
		last_imc_ban = ban
		if true {
			break
		}
	}
	return ban
}
func imc_addban(what *byte) {
	var ban *IMC_BAN
	ban = imc_newban()
	ban.Name = C.strdup(what)
}
func imc_freeban(ban *IMC_BAN) {
	for {
		if ban.Name != nil {
			libc.Free(unsafe.Pointer(ban.Name))
			ban.Name = nil
		}
		if true {
			break
		}
	}
	for {
		if ban.Prev == nil {
			first_imc_ban = ban.Next
			if first_imc_ban != nil {
				first_imc_ban.Prev = nil
			}
		} else {
			ban.Prev.Next = ban.Next
		}
		if ban.Next == nil {
			last_imc_ban = ban.Prev
			if last_imc_ban != nil {
				last_imc_ban.Next = nil
			}
		} else {
			ban.Next.Prev = ban.Prev
		}
		if true {
			break
		}
	}
	for {
		if ban != nil {
			libc.Free(unsafe.Pointer(ban))
			ban = nil
		}
		if true {
			break
		}
	}
}
func imc_delban(what *byte) bool {
	var (
		ban      *IMC_BAN
		ban_next *IMC_BAN
	)
	for ban = first_imc_ban; ban != nil; ban = ban_next {
		ban_next = ban.Next
		if C.strcasecmp(what, ban.Name) == 0 {
			imc_freeban(ban)
			return TRUE != 0
		}
	}
	return FALSE != 0
}
func imc_findchannel(name *byte) *IMC_CHANNEL {
	var c *IMC_CHANNEL
	for c = first_imc_channel; c != nil; c = c.Next {
		if c.Name != nil && C.strcasecmp(c.Name, name) == 0 || c.Local_name != nil && C.strcasecmp(c.Local_name, name) == 0 {
			return c
		}
	}
	return nil
}
func imc_freechan(c *IMC_CHANNEL) {
	var x int
	if c == nil {
		imcbug(libc.CString("%s"), "imc_freechan: Freeing NULL channel!")
		return
	}
	for {
		if c.Prev == nil {
			first_imc_channel = c.Next
			if first_imc_channel != nil {
				first_imc_channel.Prev = nil
			}
		} else {
			c.Prev.Next = c.Next
		}
		if c.Next == nil {
			last_imc_channel = c.Prev
			if last_imc_channel != nil {
				last_imc_channel.Next = nil
			}
		} else {
			c.Next.Prev = c.Prev
		}
		if true {
			break
		}
	}
	for {
		if c.Name != nil {
			libc.Free(unsafe.Pointer(c.Name))
			c.Name = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Owner != nil {
			libc.Free(unsafe.Pointer(c.Owner))
			c.Owner = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Operators != nil {
			libc.Free(unsafe.Pointer(c.Operators))
			c.Operators = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Invited != nil {
			libc.Free(unsafe.Pointer(c.Invited))
			c.Invited = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Excluded != nil {
			libc.Free(unsafe.Pointer(c.Excluded))
			c.Excluded = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Local_name != nil {
			libc.Free(unsafe.Pointer(c.Local_name))
			c.Local_name = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Regformat != nil {
			libc.Free(unsafe.Pointer(c.Regformat))
			c.Regformat = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Emoteformat != nil {
			libc.Free(unsafe.Pointer(c.Emoteformat))
			c.Emoteformat = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Socformat != nil {
			libc.Free(unsafe.Pointer(c.Socformat))
			c.Socformat = nil
		}
		if true {
			break
		}
	}
	for x = 0; x < MAX_IMCHISTORY; x++ {
		for {
			if (c.History[x]) != nil {
				libc.Free(unsafe.Pointer(c.History[x]))
				c.History[x] = nil
			}
			if true {
				break
			}
		}
	}
	for {
		if c != nil {
			libc.Free(unsafe.Pointer(c))
			c = nil
		}
		if true {
			break
		}
	}
}
func imcformat_channel(ch *char_data, d *IMC_CHANNEL, format int, all bool) {
	var (
		c   *IMC_CHANNEL = nil
		buf [4096]byte
	)
	if all {
		for c = first_imc_channel; c != nil; c = c.Next {
			if c.Local_name == nil || *c.Local_name == '\x00' {
				continue
			}
			if format == 1 || format == 4 {
				stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~C%%s: ~c%%s", c.Local_name)
				for {
					if c.Regformat != nil {
						libc.Free(unsafe.Pointer(c.Regformat))
						c.Regformat = nil
					}
					if true {
						break
					}
				}
				c.Regformat = C.strdup(&buf[0])
			}
			if format == 2 || format == 4 {
				stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~c%%s %%s", c.Local_name)
				for {
					if c.Emoteformat != nil {
						libc.Free(unsafe.Pointer(c.Emoteformat))
						c.Emoteformat = nil
					}
					if true {
						break
					}
				}
				c.Emoteformat = C.strdup(&buf[0])
			}
			if format == 3 || format == 4 {
				stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~c%%s", c.Local_name)
				for {
					if c.Socformat != nil {
						libc.Free(unsafe.Pointer(c.Socformat))
						c.Socformat = nil
					}
					if true {
						break
					}
				}
				c.Socformat = C.strdup(&buf[0])
			}
		}
	} else {
		if ch != nil && (d.Local_name == nil || *d.Local_name == '\x00') {
			imc_to_char(libc.CString("This channel is not yet locally configured.\r\n"), ch)
			return
		}
		if format == 1 || format == 4 {
			stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~C%%s: ~c%%s", d.Local_name)
			for {
				if d.Regformat != nil {
					libc.Free(unsafe.Pointer(d.Regformat))
					d.Regformat = nil
				}
				if true {
					break
				}
			}
			d.Regformat = C.strdup(&buf[0])
		}
		if format == 2 || format == 4 {
			stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~c%%s %%s", d.Local_name)
			for {
				if d.Emoteformat != nil {
					libc.Free(unsafe.Pointer(d.Emoteformat))
					d.Emoteformat = nil
				}
				if true {
					break
				}
			}
			d.Emoteformat = C.strdup(&buf[0])
		}
		if format == 3 || format == 4 {
			stdio.Snprintf(&buf[0], LGST, "~R[~Y%s~R] ~c%%s", d.Local_name)
			for {
				if d.Socformat != nil {
					libc.Free(unsafe.Pointer(d.Socformat))
					d.Socformat = nil
				}
				if true {
					break
				}
			}
			d.Socformat = C.strdup(&buf[0])
		}
	}
	imc_save_channels()
}
func imc_new_channel(chan_ *byte, owner *byte, ops *byte, invite *byte, exclude *byte, copen bool, perm int, lname *byte) {
	var c *IMC_CHANNEL
	if chan_ == nil || *chan_ == '\x00' {
		imclog(libc.CString("%s: NULL channel name received, skipping"), libc.FuncName())
		return
	}
	if C.strchr(chan_, ':') == nil {
		imclog(libc.CString("%s: Improperly formatted channel name: %s"), libc.FuncName(), chan_)
		return
	}
	for {
		if (func() *IMC_CHANNEL {
			c = new(IMC_CHANNEL)
			return c
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	c.Name = C.strdup(chan_)
	c.Owner = C.strdup(owner)
	c.Operators = C.strdup(ops)
	c.Invited = C.strdup(invite)
	c.Excluded = C.strdup(exclude)
	if lname != nil && *lname != '\x00' {
		c.Local_name = C.strdup(lname)
	} else {
		c.Local_name = imc_channel_nameof(c.Name)
	}
	c.Level = int16(perm)
	c.Refreshed = TRUE != 0
	c.Open = copen
	for {
		if first_imc_channel == nil {
			first_imc_channel = c
			last_imc_channel = c
		} else {
			last_imc_channel.Next = c
		}
		c.Next = nil
		if first_imc_channel == c {
			c.Prev = nil
		} else {
			c.Prev = last_imc_channel
		}
		last_imc_channel = c
		if true {
			break
		}
	}
	imcformat_channel(nil, c, 4, FALSE != 0)
}
func imcfread_number(fp *C.FILE) int {
	var (
		number int
		sign   bool
		c      int8
	)
	for {
		if C.feof(fp) != 0 {
			imclog(libc.CString("%s"), "imcfread_number: EOF encountered on read.")
			return 0
		}
		c = int8(getc(fp))
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(c))))) & int(uint16(int16(_ISspace)))) == 0 {
			break
		}
	}
	number = 0
	sign = FALSE != 0
	if int(c) == '+' {
		c = int8(getc(fp))
	} else if int(c) == '-' {
		sign = TRUE != 0
		c = int8(getc(fp))
	}
	if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(c))))) & int(uint16(int16(_ISdigit)))) == 0 {
		imclog(libc.CString("imcfread_number: bad format. (%c)"), c)
		return 0
	}
	for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(c))))) & int(uint16(int16(_ISdigit)))) != 0 {
		if C.feof(fp) != 0 {
			imclog(libc.CString("%s"), "imcfread_number: EOF encountered on read.")
			return number
		}
		number = number*10 + int(c) - '0'
		c = int8(getc(fp))
	}
	if sign {
		number = 0 - number
	}
	if int(c) == '|' {
		number += imcfread_number(fp)
	} else if int(c) != ' ' {
		ungetc(int(c), fp)
	}
	return number
}
func imcfread_line(fp *C.FILE) *byte {
	var (
		line  [4096]byte
		pline *byte
		c     int8
		ln    int
	)
	pline = &line[0]
	line[0] = '\x00'
	ln = 0
	for {
		if C.feof(fp) != 0 {
			imcbug(libc.CString("%s"), "imcfread_line: EOF encountered on read.")
			strlcpy(&line[0], libc.CString(""), LGST)
			return C.strdup(&line[0])
		}
		c = int8(getc(fp))
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(c))))) & int(uint16(int16(_ISspace)))) == 0 {
			break
		}
	}
	ungetc(int(c), fp)
	for {
		if C.feof(fp) != 0 {
			imcbug(libc.CString("%s"), "imcfread_line: EOF encountered on read.")
			*pline = '\x00'
			return C.strdup(&line[0])
		}
		c = int8(getc(fp))
		*func() *byte {
			p := &pline
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}() = byte(c)
		ln++
		if ln >= (int(LGST - 1)) {
			imcbug(libc.CString("%s"), "imcfread_line: line too long")
			break
		}
		if int(c) == '\n' || int(c) == '\r' {
			break
		}
	}
	for {
		c = int8(getc(fp))
		if int(c) != '\n' && int(c) != '\r' {
			break
		}
	}
	ungetc(int(c), fp)
	pline = (*byte)(unsafe.Add(unsafe.Pointer(pline), -1))
	*pline = '\x00'
	if line[C.strlen(&line[0])-1] == '~' {
		line[C.strlen(&line[0])-1] = '\x00'
	}
	return C.strdup(&line[0])
}
func imcfread_word(fp *C.FILE) *byte {
	var (
		word  [1024]byte
		pword *byte
		cEnd  int8
	)
	for {
		if C.feof(fp) != 0 {
			word[0] = '\x00'
			return &word[0]
		}
		cEnd = int8(getc(fp))
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(cEnd))))) & int(uint16(int16(_ISspace)))) == 0 {
			break
		}
	}
	if int(cEnd) == '\'' || int(cEnd) == '"' {
		pword = &word[0]
	} else {
		word[0] = byte(cEnd)
		pword = &word[1]
		cEnd = ' '
	}
	for ; uintptr(unsafe.Pointer(pword)) < uintptr(unsafe.Pointer(&word[SMST])); pword = (*byte)(unsafe.Add(unsafe.Pointer(pword), 1)) {
		if C.feof(fp) != 0 {
			*pword = '\x00'
			return &word[0]
		}
		*pword = byte(int8(getc(fp)))
		if func() int {
			if int(cEnd) == ' ' {
				return int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*pword))))) & int(uint16(int16(_ISspace)))
			}
			return int(libc.BoolToInt(*pword == byte(cEnd)))
		}() != 0 {
			if int(cEnd) == ' ' {
				ungetc(int(*pword), fp)
			}
			*pword = '\x00'
			return &word[0]
		}
	}
	imclog(libc.CString("%s: word too long"), libc.FuncName())
	return nil
}
func imcfread_to_eol(fp *C.FILE) {
	var c int8
	for {
		if C.feof(fp) != 0 {
			imclog(libc.CString("%s"), "imcfread_to_eol: EOF encountered on read.")
			return
		}
		c = int8(getc(fp))
		if int(c) == '\n' || int(c) == '\r' {
			break
		}
	}
	for {
		c = int8(getc(fp))
		if int(c) != '\n' && int(c) != '\r' {
			break
		}
	}
	ungetc(int(c), fp)
}
func imcfread_letter(fp *C.FILE) int8 {
	var c int8
	for {
		if C.feof(fp) != 0 {
			imclog(libc.CString("%s"), "imcfread_letter: EOF encountered on read.")
			return '\x00'
		}
		c = int8(getc(fp))
		if (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(c))))) & int(uint16(int16(_ISspace)))) == 0 {
			break
		}
	}
	return c
}
func imc_register_packet_handler(name *byte, func_ PACKET_FUN) {
	var ph *IMC_PHANDLER
	for ph = first_phandler; ph != nil; ph = ph.Next {
		if C.strcasecmp(ph.Name, name) == 0 {
			imclog(libc.CString("Unable to register packet type %s. Another module has already registered it."), name)
			return
		}
	}
	for {
		if (func() *IMC_PHANDLER {
			ph = new(IMC_PHANDLER)
			return ph
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	ph.Name = C.strdup(name)
	ph.Func = func_
	for {
		if first_phandler == nil {
			first_phandler = ph
			last_phandler = ph
		} else {
			last_phandler.Next = ph
		}
		ph.Next = nil
		if first_phandler == ph {
			ph.Prev = nil
		} else {
			ph.Prev = last_phandler
		}
		last_phandler = ph
		if true {
			break
		}
	}
}
func imc_freepacket(p *IMC_PACKET) {
	var (
		data      *IMC_PDATA
		data_next *IMC_PDATA
	)
	for data = p.First_data; data != nil; data = data_next {
		data_next = data.Next
		for {
			if data.Prev == nil {
				p.First_data = data.Next
				if p.First_data != nil {
					p.First_data.Prev = nil
				}
			} else {
				data.Prev.Next = data.Next
			}
			if data.Next == nil {
				p.Last_data = data.Prev
				if p.Last_data != nil {
					p.Last_data.Next = nil
				}
			} else {
				data.Next.Prev = data.Prev
			}
			if true {
				break
			}
		}
		for {
			if data != nil {
				libc.Free(unsafe.Pointer(data))
				data = nil
			}
			if true {
				break
			}
		}
	}
	for {
		if p != nil {
			libc.Free(unsafe.Pointer(p))
			p = nil
		}
		if true {
			break
		}
	}
}
func find_next_esign(string_ *byte, current int) int {
	var quote bool = FALSE != 0
	if *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) == '=' {
		current++
	}
	for ; *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) != '\x00'; current++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) == '\\' && *(*byte)(unsafe.Add(unsafe.Pointer(string_), current+1)) == '"' {
			current++
			continue
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) == '"' {
			quote = !quote
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) == '=' && !quote {
			break
		}
	}
	if *(*byte)(unsafe.Add(unsafe.Pointer(string_), current)) == '\x00' {
		return -1
	}
	return current
}
func imc_getData(output *byte, key *byte, packet *byte) *byte {
	var (
		current int    = 0
		i       uint64 = 0
		quote   bool   = FALSE != 0
	)
	*output = '\x00'
	if packet == nil || *packet == '\x00' || key == nil || *key == '\x00' {
		imcbug(libc.CString("%s: Invalid input"), libc.FuncName())
		return output
	}
	for (func() int {
		current = find_next_esign(packet, current)
		return current
	}()) >= 0 {
		if uint64(C.strlen(key)) > uint64(current) {
			continue
		}
		i = uint64(current - int(C.strlen(key)))
		if C.strncasecmp((*byte)(unsafe.Add(unsafe.Pointer(packet), i)), key, uint64(C.strlen(key))) == 0 {
			break
		}
	}
	if current < 0 {
		return output
	}
	current++
	if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == '"' {
		quote = TRUE != 0
		current++
	}
	for i = 0; *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) != '\x00'; current++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == '"' && quote {
			break
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == ' ' && !quote {
			break
		}
		if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) != '\\' {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = *(*byte)(unsafe.Add(unsafe.Pointer(packet), current))
			continue
		}
		current++
		if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == 'n' {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = '\n'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == 'r' {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = '\r'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == '"' {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = '"'
		} else if *(*byte)(unsafe.Add(unsafe.Pointer(packet), current)) == '\\' {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = '\\'
		} else {
			*(*byte)(unsafe.Add(unsafe.Pointer(output), func() uint64 {
				p := &i
				x := *p
				*p++
				return x
			}())) = *(*byte)(unsafe.Add(unsafe.Pointer(packet), current))
		}
	}
	*(*byte)(unsafe.Add(unsafe.Pointer(output), i)) = '\x00'
	return output
}
func imc_write_buffer(txt *byte) {
	var (
		output [16384]byte
		length uint64
	)
	if this_imcmud == nil || this_imcmud.Desc < 1 {
		imcbug(libc.CString("%s: Configuration or socket is invalid!"), libc.FuncName())
		return
	}
	if this_imcmud.Outbuf == nil {
		imcbug(libc.CString("%s: Output buffer has not been allocated!"), libc.FuncName())
		return
	}
	stdio.Snprintf(&output[0], IMC_BUFF_SIZE, "%s\r\n", txt)
	length = uint64(C.strlen(&output[0]))
	for this_imcmud.Outtop+int(length) >= int(this_imcmud.Outsize) {
		if this_imcmud.Outsize > 64000 {
			this_imcmud.Outtop = 0
			imcbug(libc.CString("Buffer overflow: %ld. Purging."), this_imcmud.Outsize)
			return
		}
		this_imcmud.Outsize *= 2
		for {
			if (func() *byte {
				p := &this_imcmud.Outbuf
				this_imcmud.Outbuf = (*byte)(libc.Realloc(unsafe.Pointer(this_imcmud.Outbuf), int(this_imcmud.Outsize*uint(unsafe.Sizeof(int8(0))))))
				return *p
			}()) == nil {
				imclog(libc.CString("Realloc failure @ %s:%d\n"), __FILE__, __LINE__)
				abort()
			}
			if true {
				break
			}
		}
	}
	strncpy((*byte)(unsafe.Add(unsafe.Pointer(this_imcmud.Outbuf), this_imcmud.Outtop)), &output[0], length)
	this_imcmud.Outtop += int(length)
	*(*byte)(unsafe.Add(unsafe.Pointer(this_imcmud.Outbuf), this_imcmud.Outtop)) = '\x00'
}
func imc_write_packet(p *IMC_PACKET) {
	var (
		data *IMC_PDATA
		txt  [16384]byte
	)
	stdio.Snprintf(&txt[0], IMC_BUFF_SIZE, "%s %lu %s %s %s", &p.From[0], func() uint {
		p := &imc_sequencenumber
		*p++
		return *p
	}(), this_imcmud.Localname, &p.Type[0], &p.To[0])
	for data = p.First_data; data != nil; data = data.Next {
		stdio.Snprintf(&txt[C.strlen(&txt[0])], int(IMC_BUFF_SIZE-C.strlen(&txt[0])), "%s", &data.Field[0])
	}
	imc_freepacket(p)
	imc_write_buffer(&txt[0])
}
func imc_addtopacket(p *IMC_PACKET, fmt *byte, _rest ...interface{}) {
	var (
		data *IMC_PDATA
		pkt  [16384]byte
		args libc.ArgList
	)
	args.Start(fmt, _rest)
	stdio.Vsnprintf(&pkt[0], IMC_BUFF_SIZE, libc.GoString(fmt), args)
	args.End()
	for {
		if (func() *IMC_PDATA {
			data = new(IMC_PDATA)
			return data
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	stdio.Snprintf(&data.Field[0], IMC_BUFF_SIZE, " %s", escape_string(&pkt[0]))
	for {
		if p.First_data == nil {
			p.First_data = data
			p.Last_data = data
		} else {
			p.Last_data.Next = data
		}
		data.Next = nil
		if p.First_data == data {
			data.Prev = nil
		} else {
			data.Prev = p.Last_data
		}
		p.Last_data = data
		if true {
			break
		}
	}
}
func imc_newpacket(from *byte, type_ *byte, to *byte) *IMC_PACKET {
	var p *IMC_PACKET
	if type_ == nil || *type_ == '\x00' {
		imcbug(libc.CString("%s: Attempt to build packet with no type field."), libc.FuncName())
		return nil
	}
	if from == nil || *from == '\x00' {
		imcbug(libc.CString("%s: Attempt to build %s packet with no from field."), libc.FuncName(), type_)
		return nil
	}
	if to == nil || *to == '\x00' {
		imcbug(libc.CString("%s: Attempt to build %s packet with no to field."), libc.FuncName(), type_)
		return nil
	}
	for {
		if (func() *IMC_PACKET {
			p = new(IMC_PACKET)
			return p
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	stdio.Snprintf(&p.From[0], SMST, "%s@%s", from, this_imcmud.Localname)
	strlcpy(&p.Type[0], type_, SMST)
	strlcpy(&p.To[0], to, SMST)
	p.First_data = func() *IMC_PDATA {
		p := &p.Last_data
		p.Last_data = nil
		return *p
	}()
	return p
}
func imc_update_tellhistory(ch *char_data, msg *byte) {
	var (
		new_msg [4096]byte
		local   *tm = C.localtime(&imc_time)
		x       int
	)
	stdio.Snprintf(&new_msg[0], LGST, "~R[%-2.2d:%-2.2d] %s", local.Tm_hour, local.Tm_min, msg)
	for x = 0; x < MAX_IMCTELLHISTORY; x++ {
		if uintptr(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_tellhistory[x])) == uintptr('\x00') {
			ch.Player_specials.Imcchardata.Imc_tellhistory[x] = C.strdup(&new_msg[0])
			break
		}
		if x == int(MAX_IMCTELLHISTORY-1) {
			var i int
			for i = 1; i < MAX_IMCTELLHISTORY; i++ {
				for {
					if (ch.Player_specials.Imcchardata.Imc_tellhistory[i-1]) != nil {
						libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_tellhistory[i-1]))
						ch.Player_specials.Imcchardata.Imc_tellhistory[i-1] = nil
					}
					if true {
						break
					}
				}
				ch.Player_specials.Imcchardata.Imc_tellhistory[i-1] = C.strdup(ch.Player_specials.Imcchardata.Imc_tellhistory[i])
			}
			for {
				if (ch.Player_specials.Imcchardata.Imc_tellhistory[x]) != nil {
					libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_tellhistory[x]))
					ch.Player_specials.Imcchardata.Imc_tellhistory[x] = nil
				}
				if true {
					break
				}
			}
			ch.Player_specials.Imcchardata.Imc_tellhistory[x] = C.strdup(&new_msg[0])
		}
	}
}
func imc_send_tell(from *byte, to *byte, txt *byte, reply int) {
	var p *IMC_PACKET
	p = imc_newpacket(from, libc.CString("tell"), to)
	imc_addtopacket(p, libc.CString("text=%s"), txt)
	if reply > 0 {
		imc_addtopacket(p, libc.CString("isreply=%d"), reply)
	}
	imc_write_packet(p)
}
func imc_recv_tell(q *IMC_PACKET, packet *byte) {
	var (
		vic     *char_data
		txt     [4096]byte
		isreply [1024]byte
		buf     [4096]byte
		reply   int
	)
	imc_getData(&txt[0], libc.CString("text"), packet)
	imc_getData(&isreply[0], libc.CString("isreply"), packet)
	reply = libc.Atoi(libc.GoString(&isreply[0]))
	if reply < 0 || reply > 2 {
		reply = 0
	}
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) == nil || vic.Player_specials.Imcchardata.Imcperm < IMCPERM_MORT {
		stdio.Snprintf(&buf[0], LGST, "No player named %s exists here.", &q.To[0])
		imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		return
	}
	if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("ICE")) != 0 {
		if (vic.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
			if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
				stdio.Snprintf(&buf[0], LGST, "%s is not receiving tells.", &q.To[0])
				imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
			}
			return
		}
		if imc_isignoring(vic, &q.From[0]) {
			if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
				stdio.Snprintf(&buf[0], LGST, "%s is not receiving tells.", &q.To[0])
				imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
			}
			return
		}
		if (vic.Player_specials.Imcchardata.Imcflag&(1<<0)) != 0 || (vic.Player_specials.Imcchardata.Imcflag&(1<<1)) != 0 {
			if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
				stdio.Snprintf(&buf[0], LGST, "%s is not receiving tells.", &q.To[0])
				imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
			}
			return
		}
		if (vic.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
			if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
				stdio.Snprintf(&buf[0], LGST, "%s is currently AFK. Try back later.", &q.To[0])
				imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
			}
			return
		}
		if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
			for {
				if vic.Player_specials.Imcchardata.Rreply != nil {
					libc.Free(unsafe.Pointer(vic.Player_specials.Imcchardata.Rreply))
					vic.Player_specials.Imcchardata.Rreply = nil
				}
				if true {
					break
				}
			}
			for {
				if vic.Player_specials.Imcchardata.Rreply_name != nil {
					libc.Free(unsafe.Pointer(vic.Player_specials.Imcchardata.Rreply_name))
					vic.Player_specials.Imcchardata.Rreply_name = nil
				}
				if true {
					break
				}
			}
			vic.Player_specials.Imcchardata.Rreply = C.strdup(&q.From[0])
			vic.Player_specials.Imcchardata.Rreply_name = C.strdup(imcgetname(&q.From[0]))
		}
	}
	if reply == 2 {
		stdio.Snprintf(&buf[0], LGST, "~WImctell: ~c%s\r\n", &txt[0])
	} else {
		stdio.Snprintf(&buf[0], LGST, "~C%s ~cimctells you ~c'~W%s~c'~!\r\n", imcgetname(&q.From[0]), &txt[0])
	}
	imc_to_char(&buf[0], vic)
	imc_update_tellhistory(vic, &buf[0])
}
func imc_recv_emote(q *IMC_PACKET, packet *byte) {
	var (
		d     *descriptor_data
		ch    *char_data
		txt   [4096]byte
		lvl   [1024]byte
		level int
	)
	imc_getData(&txt[0], libc.CString("text"), packet)
	imc_getData(&lvl[0], libc.CString("level"), packet)
	level = get_imcpermvalue(&lvl[0])
	if level < 0 || level > IMCPERM_IMP {
		level = IMCPERM_IMM
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Connected == CON_PLAYING && (func() *char_data {
			ch = func() *char_data {
				if d.Original != nil {
					return d.Original
				}
				return d.Character
			}()
			return ch
		}()) != nil && ch.Player_specials.Imcchardata.Imcperm >= level {
			imc_printf(ch, libc.CString("~p[~GIMC~p] %s %s\r\n"), imcgetname(&q.From[0]), &txt[0])
		}
	}
}
func update_imchistory(channel *IMC_CHANNEL, message *byte) {
	var (
		msg   [4096]byte
		buf   [4096]byte
		local *tm
		x     int
	)
	if channel == nil {
		imcbug(libc.CString("%s"), "update_imchistory: NULL channel received!")
		return
	}
	if message == nil || *message == '\x00' {
		imcbug(libc.CString("%s"), "update_imchistory: NULL message received!")
		return
	}
	strlcpy(&msg[0], message, LGST)
	for x = 0; x < MAX_IMCHISTORY; x++ {
		if channel.History[x] == nil {
			local = C.localtime(&imc_time)
			stdio.Snprintf(&buf[0], LGST, "~R[%-2.2d/%-2.2d %-2.2d:%-2.2d] ~G%s", local.Tm_mon+1, local.Tm_mday, local.Tm_hour, local.Tm_min, &msg[0])
			channel.History[x] = C.strdup(&buf[0])
			if (channel.Flags & (1 << 0)) != 0 {
				var fp *C.FILE
				stdio.Snprintf(&buf[0], LGST, "%s%s.log", IMC_DIR, channel.Local_name)
				if (func() *C.FILE {
					fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&buf[0]), "a")))
					return fp
				}()) == nil {
					C.perror(&buf[0])
					imcbug(libc.CString("Could not open file %s!"), &buf[0])
				} else {
					stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n", imc_strip_colors(channel.History[x]))
					C.fclose(fp)
					fp = nil
				}
			}
			break
		}
		if x == int(MAX_IMCHISTORY-1) {
			var y int
			for y = 1; y < MAX_IMCHISTORY; y++ {
				var z int = y - 1
				if channel.History[z] != nil {
					for {
						if (channel.History[z]) != nil {
							libc.Free(unsafe.Pointer(channel.History[z]))
							channel.History[z] = nil
						}
						if true {
							break
						}
					}
					channel.History[z] = C.strdup(channel.History[y])
				}
			}
			local = C.localtime(&imc_time)
			stdio.Snprintf(&buf[0], LGST, "~R[%-2.2d/%-2.2d %-2.2d:%-2.2d] ~G%s", local.Tm_mon+1, local.Tm_mday, local.Tm_hour, local.Tm_min, &msg[0])
			for {
				if (channel.History[x]) != nil {
					libc.Free(unsafe.Pointer(channel.History[x]))
					channel.History[x] = nil
				}
				if true {
					break
				}
			}
			channel.History[x] = C.strdup(&buf[0])
			if (channel.Flags & (1 << 0)) != 0 {
				var fp *C.FILE
				stdio.Snprintf(&buf[0], LGST, "%s%s.log", IMC_DIR, channel.Local_name)
				if (func() *C.FILE {
					fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&buf[0]), "a")))
					return fp
				}()) == nil {
					C.perror(&buf[0])
					imcbug(libc.CString("Could not open file %s!"), &buf[0])
				} else {
					stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s\n", imc_strip_colors(channel.History[x]))
					C.fclose(fp)
					fp = nil
				}
			}
		}
	}
}
func imc_display_channel(c *IMC_CHANNEL, from *byte, txt *byte, emote int) {
	var (
		d    *descriptor_data
		ch   *char_data
		buf  [4096]byte
		name [1024]byte
	)
	if c.Local_name == nil || *c.Local_name == '\x00' || !c.Refreshed {
		return
	}
	if emote < 2 {
		stdio.Snprintf(&buf[0], LGST, libc.GoString(func() *byte {
			if emote != 0 {
				return c.Emoteformat
			}
			return c.Regformat
		}()), from, txt)
	} else {
		stdio.Snprintf(&buf[0], LGST, libc.GoString(c.Socformat), txt)
	}
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Original != nil {
			ch = d.Original
		} else {
			ch = d.Character
		}
		if ch == nil || d.Connected != CON_PLAYING {
			continue
		}
		if IS_NPC(ch) {
			continue
		}
		if ch.Player_specials.Imcchardata.Imcperm < int(c.Level) || !imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
			continue
		}
		if imc_isignoring(ch, from) {
			continue
		}
		if !c.Open {
			stdio.Snprintf(&name[0], SMST, "%s@%s", GET_NAME(ch), this_imcmud.Localname)
			if !imc_hasname(c.Invited, &name[0]) && C.strcasecmp(c.Owner, &name[0]) != 0 {
				continue
			}
		}
		imc_printf(ch, libc.CString("%s\r\n"), &buf[0])
	}
	update_imchistory(c, &buf[0])
}
func imc_recv_pbroadcast(q *IMC_PACKET, packet *byte) {
	var (
		c      *IMC_CHANNEL
		chan_  [1024]byte
		txt    [4096]byte
		emote  [1024]byte
		sender [1024]byte
		em     int
	)
	imc_getData(&chan_[0], libc.CString("channel"), packet)
	imc_getData(&txt[0], libc.CString("text"), packet)
	imc_getData(&emote[0], libc.CString("emote"), packet)
	imc_getData(&sender[0], libc.CString("realfrom"), packet)
	em = libc.Atoi(libc.GoString(&emote[0]))
	if em < 0 || em > 2 {
		em = 0
	}
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		return
	}
	imc_display_channel(c, &sender[0], &txt[0], em)
}
func imc_recv_broadcast(q *IMC_PACKET, packet *byte) {
	var (
		c      *IMC_CHANNEL
		chan_  [1024]byte
		txt    [4096]byte
		emote  [1024]byte
		sender [1024]byte
		em     int
	)
	imc_getData(&chan_[0], libc.CString("channel"), packet)
	imc_getData(&txt[0], libc.CString("text"), packet)
	imc_getData(&emote[0], libc.CString("emote"), packet)
	imc_getData(&sender[0], libc.CString("sender"), packet)
	em = libc.Atoi(libc.GoString(&emote[0]))
	if em < 0 || em > 2 {
		em = 0
	}
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		return
	}
	if sender[0] == '\x00' {
		imc_display_channel(c, &q.From[0], &txt[0], em)
	} else {
		imc_display_channel(c, &sender[0], &txt[0], em)
	}
}
func imc_sendmessage(c *IMC_CHANNEL, name *byte, text *byte, emote int) {
	var p *IMC_PACKET
	if !c.Open {
		var to [1024]byte
		stdio.Snprintf(&to[0], SMST, "IMC@%s", imc_channel_mudof(c.Name))
		p = imc_newpacket(name, libc.CString("ice-msg-p"), &to[0])
	} else {
		p = imc_newpacket(name, libc.CString("ice-msg-b"), libc.CString("*@*"))
	}
	imc_addtopacket(p, libc.CString("channel=%s"), c.Name)
	imc_addtopacket(p, libc.CString("text=%s"), text)
	imc_addtopacket(p, libc.CString("emote=%d"), emote)
	imc_addtopacket(p, libc.CString("%s"), "echo=1")
	imc_write_packet(p)
}
func imc_recv_chanwhoreply(q *IMC_PACKET, packet *byte) {
	var (
		c     *IMC_CHANNEL
		vic   *char_data
		chan_ [1024]byte
		list  [16384]byte
	)
	imc_getData(&chan_[0], libc.CString("channel"), packet)
	imc_getData(&list[0], libc.CString("list"), packet)
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		return
	}
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) == nil {
		return
	}
	imc_printf(vic, libc.CString("~G%s"), &list[0])
}
func get_local_chanwho(c *IMC_CHANNEL) *byte {
	var (
		d      *descriptor_data
		person *char_data
		buf    [16384]byte
		count  int = 0
		col    int = 0
	)
	stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "The following people are listening to %s on %s:\r\n\r\n", c.Local_name, this_imcmud.Localname)
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Original != nil {
			person = d.Original
		} else {
			person = d.Character
		}
		if person == nil {
			continue
		}
		if (person.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
			continue
		}
		if !imc_hasname(person.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
			continue
		}
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "%-15s", GET_NAME(person))
		count++
		if func() int {
			p := &col
			*p++
			return *p
		}()%3 == 0 {
			col = 0
			stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "%s", "\r\n")
		}
	}
	if col != 0 {
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "%s", "\r\n")
	}
	if count == 0 {
		imcstrlcat(&buf[0], libc.CString("Nobody\r\n"), IMC_BUFF_SIZE)
	} else {
		imcstrlcat(&buf[0], libc.CString("\r\n"), IMC_BUFF_SIZE)
	}
	return &buf[0]
}
func imc_recv_chanwho(q *IMC_PACKET, packet *byte) {
	var (
		p       *IMC_PACKET
		c       *IMC_CHANNEL
		buf     [16384]byte
		lvl     [1024]byte
		channel [1024]byte
		lname   [1024]byte
		level   int
	)
	imc_getData(&lvl[0], libc.CString("level"), packet)
	level = get_imcpermvalue(&lvl[0])
	if level < 0 || level > IMCPERM_IMP {
		level = IMCPERM_ADMIN
	}
	imc_getData(&channel[0], libc.CString("channel"), packet)
	imc_getData(&lname[0], libc.CString("lname"), packet)
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&channel[0])
		return c
	}()) == nil {
		return
	}
	if c.Local_name == nil {
		stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "Channel %s is not locally configured on %s\r\n", &lname[0], this_imcmud.Localname)
	} else if int(c.Level) > level {
		stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "Channel %s is above your permission level on %s\r\n", &lname[0], this_imcmud.Localname)
	} else {
		var cwho [16384]byte
		strlcpy(&cwho[0], get_local_chanwho(c), IMC_BUFF_SIZE)
		if (C.strcasecmp(&cwho[0], libc.CString("")) == 0 || C.strcasecmp(&cwho[0], libc.CString("Nobody")) == 0) && C.strcasecmp(&q.To[0], libc.CString("*")) == 0 {
			return
		}
		strlcpy(&buf[0], &cwho[0], IMC_BUFF_SIZE)
	}
	p = imc_newpacket(libc.CString("*"), libc.CString("ice-chan-whoreply"), &q.From[0])
	imc_addtopacket(p, libc.CString("channel=%s"), c.Name)
	imc_addtopacket(p, libc.CString("list=%s"), &buf[0])
	imc_write_packet(p)
}
func imccenterline(string_ *byte, length int) *byte {
	var (
		stripped [300]byte
		outbuf   [400]byte
		amount   int
	)
	strlcpy(&stripped[0], imc_strip_colors(string_), 300)
	amount = length - int(C.strlen(&stripped[0]))
	if amount < 1 {
		amount = 1
	}
	stdio.Snprintf(&outbuf[0], 400, "%*s%s%*s", amount/2, "", string_, func() int {
		if ((amount / 2) * 2) == amount {
			return amount / 2
		}
		return (amount / 2) + 1
	}(), "")
	return &outbuf[0]
}
func imcrankbuffer(ch *char_data) *byte {
	var rbuf [1024]byte
	if ch.Player_specials.Imcchardata.Imcperm >= IMCPERM_IMM {
		strlcpy(&rbuf[0], libc.CString("~YStaff"), SMST)
		if (func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}()) != nil && (func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}())[0] != '\x00' && C.strcasecmp(libc.CString(func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}()), libc.CString("Male")) == 0 {
			stdio.Snprintf(&rbuf[0], SMST, "~B%-6s", color_mtoi(libc.CString(func() string {
				if ch.Sex == SEX_FEMALE {
					return "Female"
				}
				return "Male"
			}())))
		} else {
			stdio.Snprintf(&rbuf[0], SMST, "~M%-6s", color_mtoi(libc.CString(func() string {
				if ch.Sex == SEX_FEMALE {
					return "Female"
				}
				return "Male"
			}())))
		}
	} else {
		strlcpy(&rbuf[0], libc.CString("~BPlayer"), SMST)
		if (func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}()) != nil && (func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}())[0] != '\x00' && C.strcasecmp(libc.CString(func() string {
			if ch.Sex == SEX_FEMALE {
				return "Female"
			}
			return "Male"
		}()), libc.CString("Male")) == 0 {
			stdio.Snprintf(&rbuf[0], SMST, "~B%-6s", color_mtoi(libc.CString(func() string {
				if ch.Sex == SEX_FEMALE {
					return "Female"
				}
				return "Male"
			}())))
		} else {
			stdio.Snprintf(&rbuf[0], SMST, "~M%-6s", color_mtoi(libc.CString(func() string {
				if ch.Sex == SEX_FEMALE {
					return "Female"
				}
				return "Male"
			}())))
		}
	}
	return &rbuf[0]
}
func imc_send_whoreply(to *byte, txt *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(libc.CString("*"), libc.CString("who-reply"), to)
	imc_addtopacket(p, libc.CString("text=%s"), txt)
	imc_write_packet(p)
}
func imc_send_who(from *byte, to *byte, type_ *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(from, libc.CString("who"), to)
	imc_addtopacket(p, libc.CString("type=%s"), type_)
	imc_write_packet(p)
}
func break_newlines(argument *byte, arg_first *byte) *byte {
	var (
		cEnd  int8
		count int
	)
	count = 0
	if arg_first != nil {
		*arg_first = '\x00'
	}
	if argument == nil || *argument == '\x00' {
		return nil
	}
	for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	cEnd = '\n'
	if *argument == '\'' || *argument == '"' {
		cEnd = int8(*func() *byte {
			p := &argument
			x := *p
			*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
			return x
		}())
	}
	for *argument != '\x00' && func() int {
		p := &count
		*p++
		return *p
	}() <= math.MaxUint8 {
		if *argument == byte(cEnd) {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
			break
		}
		if arg_first != nil {
			*func() *byte {
				p := &arg_first
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &argument
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
		} else {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
	}
	if arg_first != nil {
		*arg_first = '\x00'
	}
	for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
	}
	return argument
}
func multiline_center(splitme *byte) *byte {
	var (
		newline [4096]byte
		arg     [1024]byte
	)
	newline[0] = '\x00'
	for {
		if *splitme == '\x00' {
			break
		}
		splitme = break_newlines(splitme, &arg[0])
		if C.strstr(&arg[0], libc.CString("<center>")) != nil {
			strlcpy(&arg[0], imcstrrep(&arg[0], libc.CString("<center>"), libc.CString("")), SMST)
			strlcpy(&arg[0], imccenterline(&arg[0], 78), SMST)
		}
		imcstrlcat(&newline[0], &arg[0], LGST)
		imcstrlcat(&newline[0], libc.CString("\n"), LGST)
	}
	return &newline[0]
}
func process_who_head(plrcount int) *byte {
	var (
		head   [4096]byte
		pcount [1024]byte
	)
	strlcpy(&head[0], whot.Head, LGST)
	stdio.Snprintf(&pcount[0], SMST, "%d", plrcount)
	strlcpy(&head[0], imcstrrep(&head[0], libc.CString("<%plrcount%>"), &pcount[0]), LGST)
	strlcpy(&head[0], multiline_center(&head[0]), LGST)
	return &head[0]
}
func process_who_tail(plrcount int) *byte {
	var (
		tail   [4096]byte
		pcount [1024]byte
	)
	strlcpy(&tail[0], whot.Tail, LGST)
	stdio.Snprintf(&pcount[0], SMST, "%d", plrcount)
	strlcpy(&tail[0], imcstrrep(&tail[0], libc.CString("<%plrcount%>"), &pcount[0]), LGST)
	strlcpy(&tail[0], multiline_center(&tail[0]), LGST)
	return &tail[0]
}
func process_plrline(plrrank *byte, plrflags *byte, plrname *byte, plrtitle *byte) *byte {
	var pline [4096]byte
	strlcpy(&pline[0], whot.Immline, LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charrank%>"), plrrank), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charflags%>"), plrflags), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charname%>"), plrname), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%chartitle%>"), plrtitle), LGST)
	imcstrlcat(&pline[0], libc.CString("\n"), LGST)
	return &pline[0]
}
func process_immline(plrrank *byte, plrflags *byte, plrname *byte, plrtitle *byte) *byte {
	var pline [4096]byte
	strlcpy(&pline[0], whot.Immline, LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charrank%>"), plrrank), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charflags%>"), plrflags), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%charname%>"), plrname), LGST)
	strlcpy(&pline[0], imcstrrep(&pline[0], libc.CString("<%chartitle%>"), plrtitle), LGST)
	imcstrlcat(&pline[0], libc.CString("\n"), LGST)
	return &pline[0]
}
func process_who_template(head *byte, tail *byte, plrlines *byte, immlines *byte, plrheader *byte, immheader *byte) *byte {
	var master [4096]byte
	strlcpy(&master[0], whot.Master, LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%head%>"), head), LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%tail%>"), tail), LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%plrheader%>"), plrheader), LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%immheader%>"), immheader), LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%plrline%>"), plrlines), LGST)
	strlcpy(&master[0], imcstrrep(&master[0], libc.CString("<%immline%>"), immlines), LGST)
	return &master[0]
}
func imc_assemble_who() *byte {
	var (
		person    *char_data
		d         *descriptor_data
		pcount    int  = 0
		plr       bool = FALSE != 0
		imm       bool = FALSE != 0
		plrheader [1024]byte
		immheader [1024]byte
		rank      [1024]byte
		flags     [1024]byte
		name      [1024]byte
		title     [1024]byte
		plrline   [1024]byte
		immline   [1024]byte
		plrlines  [4096]byte
		immlines  [4096]byte
		head      [4096]byte
		tail      [4096]byte
		master    [4096]byte
	)
	plrlines[0] = '\x00'
	immlines[0] = '\x00'
	plrheader[0] = '\x00'
	immheader[0] = '\x00'
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Original != nil {
			person = d.Original
		} else {
			person = d.Character
		}
		if person != nil && d.Connected == CON_PLAYING {
			if person.Player_specials.Imcchardata.Imcperm <= IMCPERM_NONE || person.Player_specials.Imcchardata.Imcperm >= IMCPERM_IMM {
				continue
			}
			if (person.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
				continue
			}
			pcount++
			if !plr {
				strlcpy(&plrheader[0], whot.Plrheader, SMST)
				plr = TRUE != 0
			}
			strlcpy(&rank[0], imcrankbuffer(person), SMST)
			if (person.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
				strlcpy(&flags[0], libc.CString("AFK"), SMST)
			} else {
				strlcpy(&flags[0], libc.CString("---"), SMST)
			}
			strlcpy(&name[0], GET_NAME(person), SMST)
			strlcpy(&title[0], color_mtoi(GET_TITLE(person)), SMST)
			strlcpy(&plrline[0], process_plrline(&rank[0], &flags[0], &name[0], &title[0]), SMST)
			imcstrlcat(&plrlines[0], &plrline[0], LGST)
		}
	}
	imm = FALSE != 0
	for d = descriptor_list; d != nil; d = d.Next {
		if d.Original != nil {
			person = d.Original
		} else {
			person = d.Character
		}
		if person != nil && d.Connected == CON_PLAYING {
			if person.Player_specials.Imcchardata.Imcperm <= IMCPERM_NONE || person.Player_specials.Imcchardata.Imcperm < IMCPERM_IMM {
				continue
			}
			if (person.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
				continue
			}
			pcount++
			if !imm {
				strlcpy(&immheader[0], whot.Immheader, SMST)
				imm = TRUE != 0
			}
			strlcpy(&rank[0], imcrankbuffer(person), SMST)
			if (person.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
				strlcpy(&flags[0], libc.CString("AFK"), SMST)
			} else {
				strlcpy(&flags[0], libc.CString("---"), SMST)
			}
			strlcpy(&name[0], GET_NAME(person), SMST)
			strlcpy(&title[0], color_mtoi(GET_TITLE(person)), SMST)
			strlcpy(&immline[0], process_immline(&rank[0], &flags[0], &name[0], &title[0]), SMST)
			imcstrlcat(&immlines[0], &immline[0], LGST)
		}
	}
	strlcpy(&head[0], process_who_head(pcount), LGST)
	strlcpy(&tail[0], process_who_tail(pcount), LGST)
	strlcpy(&master[0], process_who_template(&head[0], &tail[0], &plrlines[0], &immlines[0], &plrheader[0], &immheader[0]), LGST)
	return &master[0]
}
func imc_process_who(from *byte) {
	var whoreply [16384]byte
	strlcpy(&whoreply[0], imc_assemble_who(), IMC_BUFF_SIZE)
	imc_send_whoreply(from, &whoreply[0])
}
func imc_process_finger(from *byte, type_ *byte) {
	var (
		victim *char_data
		buf    [16384]byte
		to     [1024]byte
	)
	if type_ == nil || *type_ == '\x00' {
		return
	}
	type_ = imcone_argument(type_, &to[0])
	if (func() *char_data {
		victim = imc_find_user(type_)
		return victim
	}()) == nil {
		imc_send_whoreply(from, libc.CString("No such player is online.\r\n"))
		return
	}
	if (victim.Player_specials.Imcchardata.Imcflag&(1<<4)) != 0 || victim.Player_specials.Imcchardata.Imcperm < IMCPERM_MORT {
		imc_send_whoreply(from, libc.CString("No such player is online.\r\n"))
		return
	}
	stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "\r\n~cPlayer Profile for ~W%s~c:\r\n~W-------------------------------\r\n~cStatus: ~W%s\r\n~cPermission level: ~W%s\r\n~cListening to channels [Names may not match your mud]: ~W%s\r\n", GET_NAME(victim), func() string {
		if (victim.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
			return "AFK"
		}
		return "Lurking about"
	}(), imcperm_names[victim.Player_specials.Imcchardata.Imcperm], func() *byte {
		if victim.Player_specials.Imcchardata.Imc_listen != nil && *victim.Player_specials.Imcchardata.Imc_listen != '\x00' {
			return victim.Player_specials.Imcchardata.Imc_listen
		}
		return libc.CString("None")
	}())
	if (victim.Player_specials.Imcchardata.Imcflag & (1 << 5)) == 0 {
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~cEmail   : ~W%s\r\n~cHomepage: ~W%s\r\n~cICQ     : ~W%d\r\n~cAIM     : ~W%s\r\n~cYahoo   : ~W%s\r\n~cMSN     : ~W%s\r\n", func() *byte {
			if victim.Player_specials.Imcchardata.Email != nil && *victim.Player_specials.Imcchardata.Email != '\x00' {
				return victim.Player_specials.Imcchardata.Email
			}
			return libc.CString("None")
		}(), func() *byte {
			if victim.Player_specials.Imcchardata.Homepage != nil && *victim.Player_specials.Imcchardata.Homepage != '\x00' {
				return victim.Player_specials.Imcchardata.Homepage
			}
			return libc.CString("None")
		}(), victim.Player_specials.Imcchardata.Icq, func() *byte {
			if victim.Player_specials.Imcchardata.Aim != nil && *victim.Player_specials.Imcchardata.Aim != '\x00' {
				return victim.Player_specials.Imcchardata.Aim
			}
			return libc.CString("None")
		}(), func() *byte {
			if victim.Player_specials.Imcchardata.Yahoo != nil && *victim.Player_specials.Imcchardata.Yahoo != '\x00' {
				return victim.Player_specials.Imcchardata.Yahoo
			}
			return libc.CString("None")
		}(), func() *byte {
			if victim.Player_specials.Imcchardata.Msn != nil && *victim.Player_specials.Imcchardata.Msn != '\x00' {
				return victim.Player_specials.Imcchardata.Msn
			}
			return libc.CString("None")
		}())
	}
	stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~W%s\r\n", func() *byte {
		if victim.Player_specials.Imcchardata.Comment != nil && *victim.Player_specials.Imcchardata.Comment != '\x00' {
			return victim.Player_specials.Imcchardata.Comment
		}
		return libc.CString("")
	}())
	imc_send_whoreply(from, &buf[0])
}
func imc_recv_who(q *IMC_PACKET, packet *byte) {
	var (
		type_ [1024]byte
		buf   [16384]byte
	)
	imc_getData(&type_[0], libc.CString("type"), packet)
	if C.strcasecmp(&type_[0], libc.CString("who")) == 0 {
		imc_process_who(&q.From[0])
		return
	} else if C.strstr(&type_[0], libc.CString("finger")) != nil {
		imc_process_finger(&q.From[0], &type_[0])
		return
	} else if C.strcasecmp(&type_[0], libc.CString("info")) == 0 {
		stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "\r\n~WMUD Name    : ~c%s\r\n", this_imcmud.Localname)
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~WHost        : ~c%s\r\n", this_imcmud.Ihost)
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~WAdmin Email : ~c%s\r\n", this_imcmud.Email)
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~WWebsite     : ~c%s\r\n", this_imcmud.Www)
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~WIMC2 Version: ~c%s\r\n", this_imcmud.Versionid)
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(IMC_BUFF_SIZE-C.strlen(&buf[0])), "~WDetails     : ~c%s\r\n", this_imcmud.Details)
	} else {
		stdio.Snprintf(&buf[0], IMC_BUFF_SIZE, "%s is not a valid option. Options are: who, finger, or info.\r\n", &type_[0])
	}
	imc_send_whoreply(&q.From[0], &buf[0])
}
func imc_recv_whoreply(q *IMC_PACKET, packet *byte) {
	var (
		vic *char_data
		txt [16384]byte
	)
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) == nil {
		return
	}
	imc_getData(&txt[0], libc.CString("text"), packet)
	imc_to_pager(&txt[0], vic)
}
func imc_send_whoisreply(to *byte, data *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(libc.CString("*"), libc.CString("whois-reply"), to)
	imc_addtopacket(p, libc.CString("text=%s"), data)
	imc_write_packet(p)
}
func imc_recv_whoisreply(q *IMC_PACKET, packet *byte) {
	var (
		vic *char_data
		txt [4096]byte
	)
	imc_getData(&txt[0], libc.CString("text"), packet)
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) != nil {
		imc_to_char(&txt[0], vic)
	}
}
func imc_send_whois(from *byte, user *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(from, libc.CString("whois"), user)
	imc_write_packet(p)
}
func imc_recv_whois(q *IMC_PACKET, packet *byte) {
	var (
		vic *char_data
		buf [4096]byte
	)
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) != nil && (vic.Player_specials.Imcchardata.Imcflag&(1<<4)) == 0 {
		stdio.Snprintf(&buf[0], LGST, "~RIMC Locate: ~Y%s@%s: ~cOnline.\r\n", GET_NAME(vic), this_imcmud.Localname)
		imc_send_whoisreply(&q.From[0], &buf[0])
	}
}
func imc_recv_beep(q *IMC_PACKET, packet *byte) {
	var (
		vic *char_data = nil
		buf [4096]byte
	)
	if (func() *char_data {
		vic = imc_find_user(imc_nameof(&q.To[0]))
		return vic
	}()) == nil || vic.Player_specials.Imcchardata.Imcperm < IMCPERM_MORT {
		stdio.Snprintf(&buf[0], LGST, "No player named %s exists here.", &q.To[0])
		imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		return
	}
	if (vic.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
		if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
			stdio.Snprintf(&buf[0], LGST, "%s is not receiving beeps.", &q.To[0])
			imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		}
		return
	}
	if imc_isignoring(vic, &q.From[0]) {
		if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
			stdio.Snprintf(&buf[0], LGST, "%s is not receiving beeps.", &q.To[0])
			imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		}
		return
	}
	if (vic.Player_specials.Imcchardata.Imcflag&(1<<2)) != 0 || (vic.Player_specials.Imcchardata.Imcflag&(1<<3)) != 0 {
		if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
			stdio.Snprintf(&buf[0], LGST, "%s is not receiving beeps.", &q.To[0])
			imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		}
		return
	}
	if (vic.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
		if C.strcasecmp(imc_nameof(&q.From[0]), libc.CString("*")) != 0 {
			stdio.Snprintf(&buf[0], LGST, "%s is currently AFK. Try back later.", &q.To[0])
			imc_send_tell(libc.CString("*"), &q.From[0], &buf[0], 1)
		}
		return
	}
	imc_printf(vic, libc.CString("~c\a%s imcbeeps you.~!\r\n"), &q.From[0])
}
func imc_send_beep(from *byte, to *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(from, libc.CString("beep"), to)
	imc_write_packet(p)
}
func imc_recv_isalive(q *IMC_PACKET, packet *byte) {
	var (
		r       *REMOTEINFO
		version [1024]byte
		netname [1024]byte
		url     [1024]byte
		host    [1024]byte
		iport   [1024]byte
	)
	imc_getData(&version[0], libc.CString("versionid"), packet)
	imc_getData(&netname[0], libc.CString("networkname"), packet)
	imc_getData(&url[0], libc.CString("url"), packet)
	imc_getData(&host[0], libc.CString("host"), packet)
	imc_getData(&iport[0], libc.CString("port"), packet)
	if (func() *REMOTEINFO {
		r = imc_find_reminfo(imc_mudof(&q.From[0]))
		return r
	}()) == nil {
		imc_new_reminfo(imc_mudof(&q.From[0]), &version[0], &netname[0], &url[0], &q.Route[0])
		return
	}
	r.Expired = FALSE != 0
	if url[0] != '\x00' {
		for {
			if r.Url != nil {
				libc.Free(unsafe.Pointer(r.Url))
				r.Url = nil
			}
			if true {
				break
			}
		}
		r.Url = C.strdup(&url[0])
	}
	if version[0] != '\x00' {
		for {
			if r.Version != nil {
				libc.Free(unsafe.Pointer(r.Version))
				r.Version = nil
			}
			if true {
				break
			}
		}
		r.Version = C.strdup(&version[0])
	}
	if netname[0] != '\x00' {
		for {
			if r.Network != nil {
				libc.Free(unsafe.Pointer(r.Network))
				r.Network = nil
			}
			if true {
				break
			}
		}
		r.Network = C.strdup(&netname[0])
	}
	if q.Route != nil && q.Route[0] != '\x00' {
		for {
			if r.Path != nil {
				libc.Free(unsafe.Pointer(r.Path))
				r.Path = nil
			}
			if true {
				break
			}
		}
		r.Path = C.strdup(&q.Route[0])
	}
	if host[0] != '\x00' {
		for {
			if r.Host != nil {
				libc.Free(unsafe.Pointer(r.Host))
				r.Host = nil
			}
			if true {
				break
			}
		}
		r.Host = C.strdup(&host[0])
	}
	if iport[0] != '\x00' {
		for {
			if r.Port != nil {
				libc.Free(unsafe.Pointer(r.Port))
				r.Port = nil
			}
			if true {
				break
			}
		}
		r.Port = C.strdup(&iport[0])
	}
}
func imc_send_keepalive(q *IMC_PACKET, packet *byte) {
	var p *IMC_PACKET
	if q != nil {
		p = imc_newpacket(libc.CString("*"), libc.CString("is-alive"), &q.From[0])
	} else {
		p = imc_newpacket(libc.CString("*"), libc.CString("is-alive"), packet)
	}
	imc_addtopacket(p, libc.CString("versionid=%s"), this_imcmud.Versionid)
	imc_addtopacket(p, libc.CString("url=%s"), this_imcmud.Www)
	imc_addtopacket(p, libc.CString("host=%s"), this_imcmud.Ihost)
	imc_addtopacket(p, libc.CString("port=%d"), this_imcmud.Iport)
	imc_write_packet(p)
}
func imc_request_keepalive() {
	var p *IMC_PACKET
	p = imc_newpacket(libc.CString("*"), libc.CString("keepalive-request"), libc.CString("*@*"))
	imc_write_packet(p)
	imc_send_keepalive(nil, libc.CString("*@*"))
}
func imc_firstrefresh() {
	var p *IMC_PACKET
	p = imc_newpacket(libc.CString("*"), libc.CString("ice-refresh"), libc.CString("IMC@$"))
	imc_write_packet(p)
}
func imc_recv_iceupdate(q *IMC_PACKET, packet *byte) {
	var (
		c       *IMC_CHANNEL
		chan_   [1024]byte
		owner   [1024]byte
		ops     [1024]byte
		invite  [1024]byte
		exclude [1024]byte
		policy  [1024]byte
		level   [1024]byte
		lname   [1024]byte
		perm    int
		copen   bool
	)
	imc_getData(&chan_[0], libc.CString("channel"), packet)
	imc_getData(&owner[0], libc.CString("owner"), packet)
	imc_getData(&ops[0], libc.CString("operators"), packet)
	imc_getData(&invite[0], libc.CString("invited"), packet)
	imc_getData(&exclude[0], libc.CString("excluded"), packet)
	imc_getData(&policy[0], libc.CString("policy"), packet)
	imc_getData(&level[0], libc.CString("level"), packet)
	imc_getData(&lname[0], libc.CString("localname"), packet)
	if C.strcasecmp(&policy[0], libc.CString("open")) == 0 {
		copen = TRUE != 0
	} else {
		copen = FALSE != 0
	}
	perm = get_imcpermvalue(&level[0])
	if perm < 0 || perm > IMCPERM_IMP {
		perm = IMCPERM_ADMIN
	}
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		imc_new_channel(&chan_[0], &owner[0], &ops[0], &invite[0], &exclude[0], copen, perm, &lname[0])
		return
	}
	if chan_[0] == '\x00' {
		imclog(libc.CString("%s: NULL channel name received, skipping"), libc.FuncName())
		return
	}
	for {
		if c.Name != nil {
			libc.Free(unsafe.Pointer(c.Name))
			c.Name = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Owner != nil {
			libc.Free(unsafe.Pointer(c.Owner))
			c.Owner = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Operators != nil {
			libc.Free(unsafe.Pointer(c.Operators))
			c.Operators = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Invited != nil {
			libc.Free(unsafe.Pointer(c.Invited))
			c.Invited = nil
		}
		if true {
			break
		}
	}
	for {
		if c.Excluded != nil {
			libc.Free(unsafe.Pointer(c.Excluded))
			c.Excluded = nil
		}
		if true {
			break
		}
	}
	c.Name = C.strdup(&chan_[0])
	c.Owner = C.strdup(&owner[0])
	c.Operators = C.strdup(&ops[0])
	c.Invited = C.strdup(&invite[0])
	c.Excluded = C.strdup(&exclude[0])
	c.Open = copen
	if int(c.Level) == IMCPERM_NOTSET {
		c.Level = int16(perm)
	}
	c.Refreshed = TRUE != 0
}
func imc_recv_icedestroy(q *IMC_PACKET, packet *byte) {
	var (
		c     *IMC_CHANNEL
		chan_ [1024]byte
	)
	imc_getData(&chan_[0], libc.CString("channel"), packet)
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		return
	}
	imc_freechan(c)
	imc_save_channels()
}
func imctodikugender(gender int) int {
	var sex int = 0
	if gender == 0 {
		sex = SEX_MALE
	}
	if gender == 1 {
		sex = SEX_FEMALE
	}
	if gender > 1 {
		sex = SEX_NEUTRAL
	}
	return sex
}
func dikutoimcgender(gender int) int {
	var sex int = 0
	if gender > 2 || gender < 0 {
		sex = 2
	}
	if gender == SEX_MALE {
		sex = 0
	}
	if gender == SEX_FEMALE {
		sex = 1
	}
	return sex
}
func imc_get_ucache_gender(name *byte) int {
	var user *IMCUCACHE_DATA
	for user = first_imcucache; user != nil; user = user.Next {
		if C.strcasecmp(user.Name, name) == 0 {
			return user.Gender
		}
	}
	return -1
}
func imc_save_ucache() {
	var (
		fp   *C.FILE
		user *IMCUCACHE_DATA
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "Couldn't write to IMC2 ucache file.")
		return
	}
	for user = first_imcucache; user != nil; user = user.Next {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#UCACHE\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Name %s\n", user.Name)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Sex  %d\n", user.Gender)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Time %ld\n", int(user.Time))
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#END\n")
	C.fclose(fp)
	fp = nil
}
func imc_prune_ucache() {
	var (
		ucache      *IMCUCACHE_DATA
		next_ucache *IMCUCACHE_DATA
	)
	for ucache = first_imcucache; ucache != nil; ucache = next_ucache {
		next_ucache = ucache.Next
		if imc_time-ucache.Time >= 2592000 {
			for {
				if ucache.Name != nil {
					libc.Free(unsafe.Pointer(ucache.Name))
					ucache.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if ucache.Prev == nil {
					first_imcucache = ucache.Next
					if first_imcucache != nil {
						first_imcucache.Prev = nil
					}
				} else {
					ucache.Prev.Next = ucache.Next
				}
				if ucache.Next == nil {
					last_imcucache = ucache.Prev
					if last_imcucache != nil {
						last_imcucache.Next = nil
					}
				} else {
					ucache.Next.Prev = ucache.Prev
				}
				if true {
					break
				}
			}
			for {
				if ucache != nil {
					libc.Free(unsafe.Pointer(ucache))
					ucache = nil
				}
				if true {
					break
				}
			}
		}
	}
	imc_save_ucache()
}
func imc_ucache_update(name *byte, gender int) {
	var user *IMCUCACHE_DATA
	for user = first_imcucache; user != nil; user = user.Next {
		if C.strcasecmp(user.Name, name) == 0 {
			user.Gender = gender
			user.Time = imc_time
			return
		}
	}
	for {
		if (func() *IMCUCACHE_DATA {
			user = new(IMCUCACHE_DATA)
			return user
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	user.Name = C.strdup(name)
	user.Gender = gender
	user.Time = imc_time
	for {
		if first_imcucache == nil {
			first_imcucache = user
			last_imcucache = user
		} else {
			last_imcucache.Next = user
		}
		user.Next = nil
		if first_imcucache == user {
			user.Prev = nil
		} else {
			user.Prev = last_imcucache
		}
		last_imcucache = user
		if true {
			break
		}
	}
	imc_save_ucache()
}
func imc_send_ucache_update(visname *byte, gender int) {
	var p *IMC_PACKET
	p = imc_newpacket(visname, libc.CString("user-cache"), libc.CString("*@*"))
	imc_addtopacket(p, libc.CString("gender=%d"), gender)
	imc_write_packet(p)
}
func imc_recv_ucache(q *IMC_PACKET, packet *byte) {
	var (
		gen    [1024]byte
		sex    int
		gender int
	)
	imc_getData(&gen[0], libc.CString("gender"), packet)
	gender = libc.Atoi(libc.GoString(&gen[0]))
	sex = imc_get_ucache_gender(&q.From[0])
	if sex == gender {
		return
	}
	imc_ucache_update(&q.From[0], gender)
}
func imc_send_ucache_request(targetuser *byte) {
	var (
		p  *IMC_PACKET
		to [1024]byte
	)
	stdio.Snprintf(&to[0], SMST, "*@%s", imc_mudof(targetuser))
	p = imc_newpacket(libc.CString("*"), libc.CString("user-cache-request"), &to[0])
	imc_addtopacket(p, libc.CString("user=%s"), targetuser)
	imc_write_packet(p)
}
func imc_recv_ucache_request(q *IMC_PACKET, packet *byte) {
	var (
		p      *IMC_PACKET
		to     [1024]byte
		user   [1024]byte
		gender int
	)
	imc_getData(&user[0], libc.CString("user"), packet)
	gender = imc_get_ucache_gender(&user[0])
	if gender == -1 {
		return
	}
	stdio.Snprintf(&to[0], SMST, "*@%s", imc_mudof(&q.From[0]))
	p = imc_newpacket(libc.CString("*"), libc.CString("user-cache-reply"), &to[0])
	imc_addtopacket(p, libc.CString("user=%s"), &user[0])
	imc_addtopacket(p, libc.CString("gender=%d"), gender)
	imc_write_packet(p)
}
func imc_recv_ucache_reply(q *IMC_PACKET, packet *byte) {
	var (
		user   [1024]byte
		gen    [1024]byte
		sex    int
		gender int
	)
	imc_getData(&user[0], libc.CString("user"), packet)
	imc_getData(&gen[0], libc.CString("gender"), packet)
	gender = libc.Atoi(libc.GoString(&gen[0]))
	sex = imc_get_ucache_gender(&user[0])
	if sex == gender {
		return
	}
	imc_ucache_update(&user[0], gender)
}
func imc_recv_closenotify(q *IMC_PACKET, packet *byte) {
	var (
		r    *REMOTEINFO
		host [1024]byte
	)
	imc_getData(&host[0], libc.CString("host"), packet)
	if (func() *REMOTEINFO {
		r = imc_find_reminfo(&host[0])
		return r
	}()) == nil {
		return
	}
	r.Expired = TRUE != 0
}
func imc_register_default_packets() {
	if default_packets_registered {
		return
	}
	imc_register_packet_handler(libc.CString("keepalive-request"), func(q *IMC_PACKET, packet *byte) {
		imc_send_keepalive(q, packet)
	})
	imc_register_packet_handler(libc.CString("is-alive"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_isalive(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-update"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_iceupdate(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-msg-r"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_pbroadcast(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-msg-b"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_broadcast(q, packet)
	})
	imc_register_packet_handler(libc.CString("user-cache"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_ucache(q, packet)
	})
	imc_register_packet_handler(libc.CString("user-cache-request"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_ucache_request(q, packet)
	})
	imc_register_packet_handler(libc.CString("user-cache-reply"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_ucache_reply(q, packet)
	})
	imc_register_packet_handler(libc.CString("tell"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_tell(q, packet)
	})
	imc_register_packet_handler(libc.CString("emote"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_emote(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-destroy"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_icedestroy(q, packet)
	})
	imc_register_packet_handler(libc.CString("who"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_who(q, packet)
	})
	imc_register_packet_handler(libc.CString("who-reply"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_whoreply(q, packet)
	})
	imc_register_packet_handler(libc.CString("whois"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_whois(q, packet)
	})
	imc_register_packet_handler(libc.CString("whois-reply"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_whoisreply(q, packet)
	})
	imc_register_packet_handler(libc.CString("beep"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_beep(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-chan-who"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_chanwho(q, packet)
	})
	imc_register_packet_handler(libc.CString("ice-chan-whoreply"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_chanwhoreply(q, packet)
	})
	imc_register_packet_handler(libc.CString("close-notify"), func(q *IMC_PACKET, packet *byte) {
		imc_recv_closenotify(q, packet)
	})
	default_packets_registered = TRUE != 0
}
func pfun_lookup(type_ *byte) PACKET_FUN {
	var ph *IMC_PHANDLER
	for ph = first_phandler; ph != nil; ph = ph.Next {
		if C.strcasecmp(type_, ph.Name) == 0 {
			return ph.Func
		}
	}
	return nil
}
func imc_parse_packet(packet *byte) {
	var (
		p    *IMC_PACKET
		pfun PACKET_FUN
		arg  [1024]byte
		seq  uint
	)
	for {
		if (func() *IMC_PACKET {
			p = new(IMC_PACKET)
			return p
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	packet = imcone_argument(packet, &p.From[0])
	packet = imcone_argument(packet, &arg[0])
	seq = uint(libc.Atoi(libc.GoString(&arg[0])))
	packet = imcone_argument(packet, &p.Route[0])
	packet = imcone_argument(packet, &p.Type[0])
	packet = imcone_argument(packet, &p.To[0])
	if imc_isbanned(&p.From[0]) {
		for {
			if p != nil {
				libc.Free(unsafe.Pointer(p))
				p = nil
			}
			if true {
				break
			}
		}
		return
	}
	pfun = pfun_lookup(&p.Type[0])
	if pfun == nil {
		if imcpacketdebug {
			imclog(libc.CString("PACKET: From %s, Seq %lu, Route %s, Type %s, To %s, EXTRA %s"), &p.From[0], seq, &p.Route[0], &p.Type[0], &p.To[0], packet)
			imclog(libc.CString("No packet handler function has been defined for %s"), &p.Type[0])
		}
		for {
			if p != nil {
				libc.Free(unsafe.Pointer(p))
				p = nil
			}
			if true {
				break
			}
		}
		return
	}
	pfun(p, packet)
	if imc_find_reminfo(imc_mudof(&p.From[0])) == nil {
		imc_new_reminfo(imc_mudof(&p.From[0]), libc.CString("Unknown"), this_imcmud.Network, libc.CString("Unknown"), &p.Route[0])
	}
	for {
		if p != nil {
			libc.Free(unsafe.Pointer(p))
			p = nil
		}
		if true {
			break
		}
	}
}
func imc_finalize_connection(name *byte, netname *byte) {
	this_imcmud.State = uint16(int16(IMC_ONLINE))
	if netname != nil && *netname != '\x00' {
		for {
			if this_imcmud.Network != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Network))
				this_imcmud.Network = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Network = C.strdup(netname)
	}
	for {
		if this_imcmud.Servername != nil {
			libc.Free(unsafe.Pointer(this_imcmud.Servername))
			this_imcmud.Servername = nil
		}
		if true {
			break
		}
	}
	this_imcmud.Servername = C.strdup(name)
	imclog(libc.CString("Connected to %s. Network ID: %s"), name, func() *byte {
		if netname != nil && *netname != '\x00' {
			return netname
		}
		return libc.CString("Unknown")
	}())
	imcconnect_attempts = 0
	imc_request_keepalive()
	imc_firstrefresh()
}
func imc_handle_autosetup(source *byte, servername *byte, cmd *byte, txt *byte, encrypt *byte) {
	if C.strcasecmp(cmd, libc.CString("reject")) == 0 {
		if C.strcasecmp(txt, libc.CString("connected")) == 0 {
			imclog(libc.CString("There is already a mud named %s connected to the network."), this_imcmud.Localname)
			imc_shutdown(FALSE != 0)
			return
		}
		if C.strcasecmp(txt, libc.CString("private")) == 0 {
			imclog(libc.CString("%s is a private server. Autosetup denied."), servername)
			imc_shutdown(FALSE != 0)
			return
		}
		if C.strcasecmp(txt, libc.CString("full")) == 0 {
			imclog(libc.CString("%s has reached its connection limit. Autosetup denied."), servername)
			imc_shutdown(FALSE != 0)
			return
		}
		if C.strcasecmp(txt, libc.CString("ban")) == 0 {
			imclog(libc.CString("%s has banned your connection. Autosetup denied."), servername)
			imc_shutdown(FALSE != 0)
			return
		}
		imclog(libc.CString("%s: Invalid 'reject' response. Autosetup failed."), servername)
		imclog(libc.CString("Data received: %s %s %s %s %s"), source, servername, cmd, txt, encrypt)
		imc_shutdown(FALSE != 0)
		return
	}
	if C.strcasecmp(cmd, libc.CString("accept")) == 0 {
		imclog(libc.CString("Autosetup completed successfully."))
		if encrypt != nil && *encrypt != '\x00' && C.strcasecmp(encrypt, libc.CString("SHA256-SET")) == 0 {
			imclog(libc.CString("SHA-256 Authentication has been enabled."))
			this_imcmud.Sha256pass = TRUE != 0
			imc_save_config()
		}
		imc_finalize_connection(servername, txt)
		return
	}
	imclog(libc.CString("%s: Invalid autosetup response."), servername)
	imclog(libc.CString("Data received: %s %s %s %s %s"), source, servername, cmd, txt, encrypt)
	imc_shutdown(FALSE != 0)
}
func imc_write_socket() bool {
	var (
		ptr      *byte = this_imcmud.Outbuf
		nleft    int   = this_imcmud.Outtop
		nwritten int   = 0
	)
	if nleft <= 0 {
		return true
	}
	for nleft > 0 {
		if (func() int {
			nwritten = int(send(this_imcmud.Desc, unsafe.Pointer(ptr), uint64(nleft), 0))
			return nwritten
		}()) <= 0 {
			if nwritten == -1 && (*__errno_location()) == EAGAIN {
				var p2 *byte = this_imcmud.Outbuf
				ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), nwritten))
				for *ptr != '\x00' {
					*func() *byte {
						p := &p2
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}() = *func() *byte {
						p := &ptr
						x := *p
						*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
						return x
					}()
				}
				*p2 = '\x00'
				this_imcmud.Outtop = int(C.strlen(this_imcmud.Outbuf))
				return TRUE != 0
			}
			if nwritten < 0 {
				imclog(libc.CString("Write error on socket: %s"), C.strerror(*__errno_location()))
			} else {
				imclog(libc.CString("%s"), "Connection close detected on socket write.")
			}
			imc_shutdown(TRUE != 0)
			return FALSE != 0
		}
		nleft -= nwritten
		ptr = (*byte)(unsafe.Add(unsafe.Pointer(ptr), nwritten))
	}
	if imcpacketdebug {
		imclog(libc.CString("Packet Sent: %s"), this_imcmud.Outbuf)
		imclog(libc.CString("Bytes sent: %d"), this_imcmud.Outtop)
	}
	*this_imcmud.Outbuf = '\x00'
	this_imcmud.Outtop = 0
	return true
}
func imc_process_authentication(packet *byte) {
	var (
		command  [1024]byte
		rname    [1024]byte
		pw       [1024]byte
		version  [1024]byte
		netname  [1024]byte
		encrypt  [1024]byte
		response [4096]byte
	)
	packet = imcone_argument(packet, &command[0])
	packet = imcone_argument(packet, &rname[0])
	packet = imcone_argument(packet, &pw[0])
	packet = imcone_argument(packet, &version[0])
	packet = imcone_argument(packet, &netname[0])
	packet = imcone_argument(packet, &encrypt[0])
	if rname[0] == '\x00' {
		imclog(libc.CString("%s"), "Incomplete authentication packet. Unable to connect.")
		imc_shutdown(FALSE != 0)
		return
	}
	if C.strcasecmp(&command[0], libc.CString("SHA256-AUTH-INIT")) == 0 {
		var (
			pwd        [1024]byte
			cryptpwd   *byte
			auth_value int = 0
		)
		if pw[0] == '\x00' {
			imclog(libc.CString("SHA-256 Authentication failure: No auth_value was returned by %s."), &rname[0])
			imc_shutdown(FALSE != 0)
			return
		}
		auth_value = libc.Atoi(libc.GoString(&pw[0]))
		stdio.Snprintf(&pwd[0], SMST, "%ld%s%s", auth_value, this_imcmud.Clientpw, this_imcmud.Serverpw)
		// TODO: replace with sha256
		cryptpwd = &pwd[0]
		stdio.Snprintf(&response[0], LGST, "SHA256-AUTH-RESP %s %s version=%d", this_imcmud.Localname, cryptpwd, IMC_VERSION)
		imc_write_buffer(&response[0])
		return
	}
	if C.strcasecmp(&command[0], libc.CString("SHA256-AUTH-APPR")) == 0 {
		imclog(libc.CString("%s"), "SHA-256 Authentication completed.")
		imc_finalize_connection(&rname[0], &pw[0])
		return
	}
	if C.strcasecmp(&command[0], libc.CString("PW")) == 0 {
		if C.strcasecmp(this_imcmud.Serverpw, &pw[0]) != 0 {
			imclog(libc.CString("%s sent an improper serverpassword."), &rname[0])
			imc_shutdown(FALSE != 0)
			return
		}
		imclog(libc.CString("%s"), "Standard Authentication completed.")
		if encrypt[0] != '\x00' && C.strcasecmp(&encrypt[0], libc.CString("SHA256-SET")) == 0 {
			imclog(libc.CString("SHA-256 Authentication has been enabled."))
			this_imcmud.Sha256pass = TRUE != 0
			imc_save_config()
		}
		imc_finalize_connection(&rname[0], &netname[0])
		return
	}
	if C.strcasecmp(&command[0], libc.CString("autosetup")) == 0 {
		imc_handle_autosetup(&command[0], &rname[0], &pw[0], &version[0], &netname[0])
		return
	}
	imclog(libc.CString("Invalid authentication response received from %s!!"), &rname[0])
	imclog(libc.CString("Data received: %s %s %s %s %s"), &command[0], &rname[0], &pw[0], &version[0], &netname[0])
	imc_shutdown(FALSE != 0)
}
func imc_read_buffer() bool {
	var (
		i     uint64 = 0
		j     uint64 = 0
		ended uint8  = 0
		k     int    = 0
	)
	if this_imcmud.Inbuf[0] == '\x00' {
		return false
	}
	k = int(C.strlen(&this_imcmud.Incomm[0]))
	if k < 0 {
		k = 0
	}
	for i = 0; this_imcmud.Inbuf[i] != '\x00' && this_imcmud.Inbuf[i] != '\n' && this_imcmud.Inbuf[i] != '\r' && i < IMC_BUFF_SIZE; i++ {
		this_imcmud.Incomm[func() int {
			p := &k
			x := *p
			*p++
			return x
		}()] = this_imcmud.Inbuf[i]
	}
	for this_imcmud.Inbuf[i] == '\n' || this_imcmud.Inbuf[i] == '\r' {
		ended = 1
		i++
	}
	this_imcmud.Incomm[k] = '\x00'
	for (func() byte {
		p := &this_imcmud.Inbuf[j]
		this_imcmud.Inbuf[j] = this_imcmud.Inbuf[i+j]
		return *p
	}()) != '\x00' {
		j++
	}
	this_imcmud.Inbuf[j] = '\x00'
	return int(ended) != 0
}
func imc_read_socket() bool {
	var (
		iStart     uint64
		iErr       uint64
		loop_count int16 = 0
		begin      bool  = TRUE != 0
	)
	iStart = uint64(C.strlen(&this_imcmud.Inbuf[0]))
	for {
		var nRead int
		if int(func() int16 {
			p := &loop_count
			*p++
			return *p
		}()) > 100 {
			break
		}
		nRead = int(recv(this_imcmud.Desc, unsafe.Pointer(&this_imcmud.Inbuf[iStart]), uint64(16384-10-uintptr(iStart)), 0))
		iErr = uint64(*__errno_location())
		if nRead > 0 {
			iStart += uint64(nRead)
			if iStart >= uint64(16384-10) {
				break
			}
			begin = FALSE != 0
		} else if nRead == 0 && int(this_imcmud.State) == IMC_ONLINE {
			if !begin {
				break
			}
			imclog(libc.CString("%s"), "Connection close detected on read of IMC2 socket.")
			return FALSE != 0
		} else if iErr == EAGAIN || iErr == EAGAIN {
			break
		} else if nRead == -1 {
			imclog(libc.CString("%s: Descriptor error on #%d: %s"), libc.FuncName(), this_imcmud.Desc, C.strerror(int(iErr)))
			return FALSE != 0
		}
	}
	this_imcmud.Inbuf[iStart] = '\x00'
	return TRUE != 0
}
func imc_loop() {
	var (
		in_set    fd_set
		out_set   fd_set
		last_time timeval
		null_time timeval
	)
	gettimeofday(&last_time, nil)
	imc_time = int64(last_time.Tv_sec)
	if imcwait > 0 {
		imcwait--
	}
	if imcwait == 1 {
		if func() int {
			p := &imcconnect_attempts
			*p++
			return *p
		}() > 3 {
			if this_imcmud.Sha256pass {
				imclog(libc.CString("%s"), "Unable to reconnect using SHA-256, trying standard authentication.")
				this_imcmud.Sha256pass = FALSE != 0
				imc_save_config()
				imcconnect_attempts = 0
			} else {
				imcwait = -2
				imclog(libc.CString("%s"), "Unable to reestablish connection to server. Abandoning reconnect.")
				return
			}
		}
		imc_startup(TRUE != 0, -1, FALSE != 0)
		return
	}
	if int(this_imcmud.State) == IMC_OFFLINE || this_imcmud.Desc == -1 {
		return
	}
	if imcucache_clock <= imc_time {
		imcucache_clock = imc_time + 86400
		imc_prune_ucache()
	}
	for {
		{
			var (
				__i   uint
				__arr *fd_set = (&in_set)
			)
			for __i = 0; __i < uint(unsafe.Sizeof(fd_set{})/unsafe.Sizeof(__fd_mask(0))); __i++ {
				__arr.__fds_bits[__i] = 0
			}
		}
		if true {
			break
		}
	}
	for {
		{
			var (
				__i   uint
				__arr *fd_set = (&out_set)
			)
			for __i = 0; __i < uint(unsafe.Sizeof(fd_set{})/unsafe.Sizeof(__fd_mask(0))); __i++ {
				__arr.__fds_bits[__i] = 0
			}
		}
		if true {
			break
		}
	}
	in_set.__fds_bits[this_imcmud.Desc/(8*int(unsafe.Sizeof(__fd_mask(0))))] |= __fd_mask(1 << (this_imcmud.Desc % (8 * int(unsafe.Sizeof(__fd_mask(0))))))
	out_set.__fds_bits[this_imcmud.Desc/(8*int(unsafe.Sizeof(__fd_mask(0))))] |= __fd_mask(1 << (this_imcmud.Desc % (8 * int(unsafe.Sizeof(__fd_mask(0))))))
	null_time.Tv_sec = __time_t(func() __suseconds_t {
		p := &null_time.Tv_usec
		null_time.Tv_usec = 0
		return *p
	}())
	if netpoll_select(this_imcmud.Desc+1, &in_set, &out_set, nil, &null_time) < 0 {
		C.perror(libc.CString("imc_loop: select: poll"))
		imc_shutdown(TRUE != 0)
		return
	}
	if (in_set.__fds_bits[this_imcmud.Desc/(8*int(unsafe.Sizeof(__fd_mask(0))))] & (__fd_mask(1 << (this_imcmud.Desc % (8 * int(unsafe.Sizeof(__fd_mask(0)))))))) != 0 {
		if !imc_read_socket() {
			if this_imcmud.Inbuf != nil && this_imcmud.Inbuf[0] != '\x00' {
				if imc_read_buffer() {
					if C.strcasecmp(&this_imcmud.Incomm[0], libc.CString("SHA-256 authentication is required.")) == 0 {
						imclog(libc.CString("%s"), "Unable to reconnect using standard authentication, trying SHA-256.")
						this_imcmud.Sha256pass = TRUE != 0
						imc_save_config()
					} else {
						imclog(libc.CString("Buffer contents: %s"), &this_imcmud.Incomm[0])
					}
				}
			}
			out_set.__fds_bits[this_imcmud.Desc/(8*int(unsafe.Sizeof(__fd_mask(0))))] &= ^(__fd_mask(1 << (this_imcmud.Desc % (8 * int(unsafe.Sizeof(__fd_mask(0)))))))
			imc_shutdown(TRUE != 0)
			return
		}
		for imc_read_buffer() {
			if imcpacketdebug {
				imclog(libc.CString("Packet received: %s"), &this_imcmud.Incomm[0])
			}
			switch this_imcmud.State {
			default:
				fallthrough
			case IMC_OFFLINE:
				fallthrough
			case IMC_AUTH1:
			case IMC_AUTH2:
				imc_process_authentication(&this_imcmud.Incomm[0])
				this_imcmud.Incomm[0] = '\x00'
			case IMC_ONLINE:
				imc_parse_packet(&this_imcmud.Incomm[0])
				this_imcmud.Incomm[0] = '\x00'
			}
		}
	}
	if this_imcmud.Desc > 0 && this_imcmud.Outtop > 0 && (out_set.__fds_bits[this_imcmud.Desc/(8*int(unsafe.Sizeof(__fd_mask(0))))]&(__fd_mask(1<<(this_imcmud.Desc%(8*int(unsafe.Sizeof(__fd_mask(0)))))))) != 0 && !imc_write_socket() {
		this_imcmud.Outtop = 0
		imc_shutdown(TRUE != 0)
	}
}
func imc_adjust_perms(ch *char_data) {
	if this_imcmud == nil {
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 9)) == 0 {
		if ch.Admlevel < this_imcmud.Minlevel {
			ch.Player_specials.Imcchardata.Imcperm = IMCPERM_NONE
		} else if ch.Admlevel >= this_imcmud.Minlevel && ch.Admlevel < this_imcmud.Immlevel {
			ch.Player_specials.Imcchardata.Imcperm = IMCPERM_MORT
		} else if ch.Admlevel >= this_imcmud.Immlevel && ch.Admlevel < this_imcmud.Adminlevel {
			ch.Player_specials.Imcchardata.Imcperm = IMCPERM_IMM
		} else if ch.Admlevel >= this_imcmud.Adminlevel && ch.Admlevel < this_imcmud.Implevel {
			ch.Player_specials.Imcchardata.Imcperm = IMCPERM_ADMIN
		} else if ch.Admlevel >= this_imcmud.Implevel {
			ch.Player_specials.Imcchardata.Imcperm = IMCPERM_IMP
		}
	}
}
func imc_char_login(ch *char_data) {
	var (
		buf    [1024]byte
		gender int
		sex    int
	)
	if this_imcmud == nil {
		return
	}
	imc_adjust_perms(ch)
	if int(this_imcmud.State) != IMC_ONLINE {
		if ch.Player_specials.Imcchardata.Imcperm >= IMCPERM_IMM && imcwait == -2 {
			imc_to_char(libc.CString("~RThe IMC2 connection is down. Attempts to reconnect were abandoned due to excessive failures.\r\n"), ch)
		}
		return
	}
	if ch.Player_specials.Imcchardata.Imcperm < IMCPERM_MORT {
		return
	}
	stdio.Snprintf(&buf[0], SMST, "%s@%s", GET_NAME(ch), this_imcmud.Localname)
	gender = imc_get_ucache_gender(&buf[0])
	sex = dikutoimcgender(int(ch.Sex))
	if gender == sex {
		return
	}
	imc_ucache_update(&buf[0], sex)
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 4)) == 0 {
		imc_send_ucache_update(GET_NAME(ch), sex)
	}
}
func imc_loadchar(ch *char_data, fp *C.FILE, word *byte) bool {
	var fMatch bool = FALSE != 0
	if IS_NPC(ch) {
		return FALSE != 0
	}
	if ch.Player_specials.Imcchardata.Imcperm == IMCPERM_NOTSET {
		imc_adjust_perms(ch)
	}
	switch *word {
	case 'I':
		if C.strcasecmp(word, libc.CString("IMCPerm")) == 0 {
			ch.Player_specials.Imcchardata.Imcperm = imcfread_number(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCEmail")) == 0 {
			ch.Player_specials.Imcchardata.Email = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCAIM")) == 0 {
			ch.Player_specials.Imcchardata.Aim = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCICQ")) == 0 {
			ch.Player_specials.Imcchardata.Icq = imcfread_number(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCYahoo")) == 0 {
			ch.Player_specials.Imcchardata.Yahoo = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCMSN")) == 0 {
			ch.Player_specials.Imcchardata.Msn = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCHomepage")) == 0 {
			ch.Player_specials.Imcchardata.Homepage = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCComment")) == 0 {
			ch.Player_specials.Imcchardata.Comment = imcfread_line(fp)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCFlags")) == 0 {
			ch.Player_specials.Imcchardata.Imcflag = imcfread_number(fp)
			imc_char_login(ch)
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMClisten")) == 0 {
			ch.Player_specials.Imcchardata.Imc_listen = imcfread_line(fp)
			if ch.Player_specials.Imcchardata.Imc_listen != nil && int(this_imcmud.State) == IMC_ONLINE {
				var (
					channel  *IMC_CHANNEL = nil
					channels *byte        = ch.Player_specials.Imcchardata.Imc_listen
					arg      [1024]byte
				)
				for {
					if *channels == '\x00' {
						break
					}
					channels = imcone_argument(channels, &arg[0])
					if (func() *IMC_CHANNEL {
						channel = imc_findchannel(&arg[0])
						return channel
					}()) == nil {
						imc_removename(&ch.Player_specials.Imcchardata.Imc_listen, &arg[0])
					}
					if channel != nil && ch.Player_specials.Imcchardata.Imcperm < int(channel.Level) {
						imc_removename(&ch.Player_specials.Imcchardata.Imc_listen, &arg[0])
					}
				}
			}
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCdeny")) == 0 {
			ch.Player_specials.Imcchardata.Imc_denied = imcfread_line(fp)
			if ch.Player_specials.Imcchardata.Imc_denied != nil && int(this_imcmud.State) == IMC_ONLINE {
				var (
					channel  *IMC_CHANNEL = nil
					channels *byte        = ch.Player_specials.Imcchardata.Imc_denied
					arg      [1024]byte
				)
				for {
					if *channels == '\x00' {
						break
					}
					channels = imcone_argument(channels, &arg[0])
					if (func() *IMC_CHANNEL {
						channel = imc_findchannel(&arg[0])
						return channel
					}()) == nil {
						imc_removename(&ch.Player_specials.Imcchardata.Imc_denied, &arg[0])
					}
					if channel != nil && ch.Player_specials.Imcchardata.Imcperm < int(channel.Level) {
						imc_removename(&ch.Player_specials.Imcchardata.Imc_denied, &arg[0])
					}
				}
			}
			fMatch = TRUE != 0
			break
		}
		if C.strcasecmp(word, libc.CString("IMCignore")) == 0 {
			var temp *IMC_IGNORE
			for {
				if (func() *IMC_IGNORE {
					temp = new(IMC_IGNORE)
					return temp
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			temp.Name = imcfread_line(fp)
			for {
				if ch.Player_specials.Imcchardata.Imcfirst_ignore == nil {
					ch.Player_specials.Imcchardata.Imcfirst_ignore = temp
					ch.Player_specials.Imcchardata.Imclast_ignore = temp
				} else {
					ch.Player_specials.Imcchardata.Imclast_ignore.Next = temp
				}
				temp.Next = nil
				if ch.Player_specials.Imcchardata.Imcfirst_ignore == temp {
					temp.Prev = nil
				} else {
					temp.Prev = ch.Player_specials.Imcchardata.Imclast_ignore
				}
				ch.Player_specials.Imcchardata.Imclast_ignore = temp
				if true {
					break
				}
			}
			fMatch = TRUE != 0
			break
		}
	}
	return fMatch
}
func imc_savechar(ch *char_data, fp *C.FILE) {
	var (
		temp *IMC_IGNORE
		last *IMC_IGNORE = nil
	)
	if IS_NPC(ch) {
		return
	}
	if ch.Player_specials.Imcchardata == nil {
		return
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCPerm      %d\n", ch.Player_specials.Imcchardata.Imcperm)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCFlags     %ld\n", ch.Player_specials.Imcchardata.Imcflag)
	if ch.Player_specials.Imcchardata.Imc_listen != nil && *ch.Player_specials.Imcchardata.Imc_listen != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCListen    %s\n", ch.Player_specials.Imcchardata.Imc_listen)
	}
	if ch.Player_specials.Imcchardata.Imc_denied != nil && *ch.Player_specials.Imcchardata.Imc_denied != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCDeny      %s\n", ch.Player_specials.Imcchardata.Imc_denied)
	}
	if ch.Player_specials.Imcchardata.Email != nil && *ch.Player_specials.Imcchardata.Email != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCEmail     %s\n", ch.Player_specials.Imcchardata.Email)
	}
	if ch.Player_specials.Imcchardata.Homepage != nil && *ch.Player_specials.Imcchardata.Homepage != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCHomepage  %s\n", ch.Player_specials.Imcchardata.Homepage)
	}
	if ch.Player_specials.Imcchardata.Icq != 0 {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCICQ       %d\n", ch.Player_specials.Imcchardata.Icq)
	}
	if ch.Player_specials.Imcchardata.Aim != nil && *ch.Player_specials.Imcchardata.Aim != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCAIM       %s\n", ch.Player_specials.Imcchardata.Aim)
	}
	if ch.Player_specials.Imcchardata.Yahoo != nil && *ch.Player_specials.Imcchardata.Yahoo != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCYahoo     %s\n", ch.Player_specials.Imcchardata.Yahoo)
	}
	if ch.Player_specials.Imcchardata.Msn != nil && *ch.Player_specials.Imcchardata.Msn != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCMSN       %s\n", ch.Player_specials.Imcchardata.Msn)
	}
	if ch.Player_specials.Imcchardata.Comment != nil && *ch.Player_specials.Imcchardata.Comment != '\x00' {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCComment   %s\n", ch.Player_specials.Imcchardata.Comment)
	}
	for temp = ch.Player_specials.Imcchardata.Imcfirst_ignore; temp != nil; temp = temp.Next {
		if last != nil {
			continue
		}
		if temp == ch.Player_specials.Imcchardata.Imclast_ignore {
			last = temp
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCignore    %s\n", temp.Name)
	}
}
func imc_freechardata(ch *char_data) {
	var (
		ign      *IMC_IGNORE
		ign_next *IMC_IGNORE
		x        int
	)
	if IS_NPC(ch) {
		return
	}
	if ch.Player_specials.Imcchardata == nil {
		return
	}
	for ign = ch.Player_specials.Imcchardata.Imcfirst_ignore; ign != nil; ign = ign_next {
		ign_next = ign.Next
		for {
			if ign.Name != nil {
				libc.Free(unsafe.Pointer(ign.Name))
				ign.Name = nil
			}
			if true {
				break
			}
		}
		for {
			if ign.Prev == nil {
				ch.Player_specials.Imcchardata.Imcfirst_ignore = ign.Next
				if ch.Player_specials.Imcchardata.Imcfirst_ignore != nil {
					ch.Player_specials.Imcchardata.Imcfirst_ignore.Prev = nil
				}
			} else {
				ign.Prev.Next = ign.Next
			}
			if ign.Next == nil {
				ch.Player_specials.Imcchardata.Imclast_ignore = ign.Prev
				if ch.Player_specials.Imcchardata.Imclast_ignore != nil {
					ch.Player_specials.Imcchardata.Imclast_ignore.Next = nil
				}
			} else {
				ign.Next.Prev = ign.Prev
			}
			if true {
				break
			}
		}
		for {
			if ign != nil {
				libc.Free(unsafe.Pointer(ign))
				ign = nil
			}
			if true {
				break
			}
		}
	}
	for x = 0; x < MAX_IMCTELLHISTORY; x++ {
		for {
			if (ch.Player_specials.Imcchardata.Imc_tellhistory[x]) != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_tellhistory[x]))
				ch.Player_specials.Imcchardata.Imc_tellhistory[x] = nil
			}
			if true {
				break
			}
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Imc_listen != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_listen))
			ch.Player_specials.Imcchardata.Imc_listen = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Imc_denied != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Imc_denied))
			ch.Player_specials.Imcchardata.Imc_denied = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Rreply != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Rreply))
			ch.Player_specials.Imcchardata.Rreply = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Rreply_name != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Rreply_name))
			ch.Player_specials.Imcchardata.Rreply_name = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Email != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Email))
			ch.Player_specials.Imcchardata.Email = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Homepage != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Homepage))
			ch.Player_specials.Imcchardata.Homepage = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Aim != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Aim))
			ch.Player_specials.Imcchardata.Aim = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Yahoo != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Yahoo))
			ch.Player_specials.Imcchardata.Yahoo = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Msn != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Msn))
			ch.Player_specials.Imcchardata.Msn = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata.Comment != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Comment))
			ch.Player_specials.Imcchardata.Comment = nil
		}
		if true {
			break
		}
	}
	for {
		if ch.Player_specials.Imcchardata != nil {
			libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata))
			ch.Player_specials.Imcchardata = nil
		}
		if true {
			break
		}
	}
}
func imc_initchar(ch *char_data) {
	if IS_NPC(ch) {
		return
	}
	for {
		if (func() *imcchar_data {
			p := &ch.Player_specials.Imcchardata
			ch.Player_specials.Imcchardata = (*imcchar_data)(unsafe.Pointer(new(IMC_CHARDATA)))
			return *p
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	ch.Player_specials.Imcchardata.Imc_listen = nil
	ch.Player_specials.Imcchardata.Imc_denied = nil
	ch.Player_specials.Imcchardata.Rreply = nil
	ch.Player_specials.Imcchardata.Rreply_name = nil
	ch.Player_specials.Imcchardata.Email = nil
	ch.Player_specials.Imcchardata.Homepage = nil
	ch.Player_specials.Imcchardata.Aim = nil
	ch.Player_specials.Imcchardata.Yahoo = nil
	ch.Player_specials.Imcchardata.Msn = nil
	ch.Player_specials.Imcchardata.Comment = nil
	ch.Player_specials.Imcchardata.Imcflag = 0
	ch.Player_specials.Imcchardata.Imcflag |= 1 << 8
	ch.Player_specials.Imcchardata.Imcfirst_ignore = nil
	ch.Player_specials.Imcchardata.Imclast_ignore = nil
	ch.Player_specials.Imcchardata.Imcperm = IMCPERM_NOTSET
}
func imc_loadhistory() {
	var (
		filename [256]byte
		tempfile *C.FILE
		tempchan *IMC_CHANNEL = nil
		x        int
	)
	for tempchan = first_imc_channel; tempchan != nil; tempchan = tempchan.Next {
		if tempchan.Local_name == nil {
			continue
		}
		stdio.Snprintf(&filename[0], 256, "%s%s.hist", IMC_DIR, tempchan.Local_name)
		if (func() *C.FILE {
			tempfile = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&filename[0]), "r")))
			return tempfile
		}()) == nil {
			continue
		}
		for x = 0; x < MAX_IMCHISTORY; x++ {
			if C.feof(tempfile) != 0 {
				tempchan.History[x] = nil
			} else {
				tempchan.History[x] = imcfread_line(tempfile)
			}
		}
		C.fclose(tempfile)
		tempfile = nil
		unlink(&filename[0])
	}
}
func imc_savehistory() {
	var (
		filename [256]byte
		tempfile *C.FILE
		tempchan *IMC_CHANNEL = nil
		x        int
	)
	for tempchan = first_imc_channel; tempchan != nil; tempchan = tempchan.Next {
		if tempchan.Local_name == nil {
			continue
		}
		if tempchan.History[0] == nil {
			continue
		}
		stdio.Snprintf(&filename[0], 256, "%s%s.hist", IMC_DIR, tempchan.Local_name)
		if (func() *C.FILE {
			tempfile = (*C.FILE)(unsafe.Pointer(stdio.FOpen(libc.GoString(&filename[0]), "w")))
			return tempfile
		}()) == nil {
			continue
		}
		for x = 0; x < MAX_IMCHISTORY; x++ {
			if tempchan.History[x] != nil {
				stdio.Fprintf((*stdio.File)(unsafe.Pointer(tempfile)), "%s\n", tempchan.History[x])
			}
		}
		C.fclose(tempfile)
		tempfile = nil
	}
}
func imc_save_channels() {
	var (
		c  *IMC_CHANNEL
		fp *C.FILE
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imcbug(libc.CString("Can't write to %s"), IMC_DIR)
		return
	}
	for c = first_imc_channel; c != nil; c = c.Next {
		if c.Local_name == nil || *c.Local_name == '\x00' {
			continue
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#IMCCHAN\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanName   %s\n", c.Name)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanLocal  %s\n", c.Local_name)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanRegF   %s\n", c.Regformat)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanEmoF   %s\n", c.Emoteformat)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanSocF   %s\n", c.Socformat)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ChanLevel  %d\n", c.Level)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#END\n")
	C.fclose(fp)
	fp = nil
}
func imc_readchannel(channel *IMC_CHANNEL, fp *C.FILE) {
	var (
		word   *byte
		fMatch bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fp) != 0 {
				return "End"
			}
			return libc.GoString(imcfread_word(fp))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '*':
			fMatch = TRUE != 0
			imcfread_to_eol(fp)
		case 'C':
			if C.strcasecmp(word, libc.CString("ChanName")) == 0 {
				channel.Name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ChanLocal")) == 0 {
				channel.Local_name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ChanRegF")) == 0 {
				channel.Regformat = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ChanEmoF")) == 0 {
				channel.Emoteformat = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ChanSocF")) == 0 {
				channel.Socformat = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ChanLevel")) == 0 {
				channel.Level = int16(imcfread_number(fp))
				fMatch = TRUE != 0
				break
			}
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				if int(channel.Level) > IMCPERM_IMP {
					if int(channel.Level) < this_imcmud.Minlevel {
						channel.Level = int16(IMCPERM_NONE)
					} else if int(channel.Level) >= this_imcmud.Minlevel && int(channel.Level) < this_imcmud.Immlevel {
						channel.Level = int16(IMCPERM_MORT)
					} else if int(channel.Level) >= this_imcmud.Immlevel && int(channel.Level) < this_imcmud.Adminlevel {
						channel.Level = int16(IMCPERM_IMM)
					} else if int(channel.Level) >= this_imcmud.Adminlevel && int(channel.Level) < this_imcmud.Implevel {
						channel.Level = int16(IMCPERM_ADMIN)
					} else if int(channel.Level) >= this_imcmud.Implevel {
						channel.Level = int16(IMCPERM_IMP)
					}
				}
			}
			return
		}
		if !fMatch {
			imcbug(libc.CString("imc_readchannel: no match: %s"), word)
		}
	}
}
func imc_loadchannels() {
	var (
		fp      *C.FILE
		channel *IMC_CHANNEL
	)
	first_imc_channel = nil
	last_imc_channel = nil
	imclog(libc.CString("%s"), "Loading channels...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imcbug(libc.CString("%s"), "Can't open imc channel file")
		return
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fp)
		if int(letter) == '*' {
			imcfread_to_eol(fp)
			continue
		}
		if int(letter) != '#' {
			imcbug(libc.CString("%s"), "imc_loadchannels: # not found.")
			break
		}
		word = imcfread_word(fp)
		if C.strcasecmp(word, libc.CString("IMCCHAN")) == 0 {
			var x int
			for {
				if (func() *IMC_CHANNEL {
					channel = new(IMC_CHANNEL)
					return channel
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			imc_readchannel(channel, fp)
			for x = 0; x < MAX_IMCHISTORY; x++ {
				channel.History[x] = nil
			}
			channel.Refreshed = FALSE != 0
			for {
				if first_imc_channel == nil {
					first_imc_channel = channel
					last_imc_channel = channel
				} else {
					last_imc_channel.Next = channel
				}
				channel.Next = nil
				if first_imc_channel == channel {
					channel.Prev = nil
				} else {
					channel.Prev = last_imc_channel
				}
				last_imc_channel = channel
				if true {
					break
				}
			}
			imclog(libc.CString("configured %s as %s"), channel.Name, channel.Local_name)
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_loadchannels: bad section: %s."), word)
			continue
		}
	}
	C.fclose(fp)
	fp = nil
}
func imc_savebans() {
	var (
		out *C.FILE
		ban *IMC_BAN
	)
	if (func() *C.FILE {
		out = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return out
	}()) == nil {
		imcbug(libc.CString("%s"), "imc_savebans: error opening ban file for write")
		return
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(out)), "%s", "#IGNORES\n")
	for ban = first_imc_ban; ban != nil; ban = ban.Next {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(out)), "%s\n", ban.Name)
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(out)), "%s", "#END\n")
	C.fclose(out)
	out = nil
}
func imc_readbans() {
	var (
		inf  *C.FILE
		word *byte
		temp [1024]byte
	)
	imclog(libc.CString("%s"), "Loading ban list...")
	if (func() *C.FILE {
		inf = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return inf
	}()) == nil {
		imcbug(libc.CString("%s"), "imc_readbans: couldn't open ban file")
		return
	}
	word = imcfread_word(inf)
	if C.strcasecmp(word, libc.CString("#IGNORES")) != 0 {
		imcbug(libc.CString("%s"), "imc_readbans: Corrupt file")
		C.fclose(inf)
		inf = nil
		return
	}
	for C.feof(inf) == 0 && ferror(inf) == 0 {
		strlcpy(&temp[0], imcfread_word(inf), SMST)
		if C.strcasecmp(&temp[0], libc.CString("#END")) == 0 {
			C.fclose(inf)
			inf = nil
			return
		}
		imc_addban(&temp[0])
	}
	if ferror(inf) != 0 {
		C.perror(libc.CString("imc_readbans"))
		C.fclose(inf)
		inf = nil
		return
	}
	C.fclose(inf)
	inf = nil
}
func imc_savecolor() {
	var (
		fp    *C.FILE
		color *IMC_COLOR
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "Couldn't write to IMC2 color file.")
		return
	}
	for color = first_imc_color; color != nil; color = color.Next {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#COLOR\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Name   %s\n", color.Name)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Mudtag %s\n", color.Mudtag)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "IMCtag %s\n", color.Imctag)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#END\n")
	C.fclose(fp)
	fp = nil
}
func imc_readcolor(color *IMC_COLOR, fp *C.FILE) {
	var (
		word   *byte
		fMatch bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fp) != 0 {
				return "End"
			}
			return libc.GoString(imcfread_word(fp))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '*':
			fMatch = TRUE != 0
			imcfread_to_eol(fp)
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				return
			}
		case 'I':
			if C.strcasecmp(word, libc.CString("IMCtag")) == 0 {
				color.Imctag = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		case 'M':
			if C.strcasecmp(word, libc.CString("Mudtag")) == 0 {
				color.Mudtag = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		case 'N':
			if C.strcasecmp(word, libc.CString("Name")) == 0 {
				color.Name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		}
		if !fMatch {
			imcbug(libc.CString("imc_readcolor: no match: %s"), word)
		}
	}
}
func imc_load_color_table() {
	var (
		fp    *C.FILE
		color *IMC_COLOR
	)
	first_imc_color = func() *IMC_COLOR {
		last_imc_color = nil
		return last_imc_color
	}()
	imclog(libc.CString("%s"), "Loading IMC2 color table...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "No color table found.")
		return
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fp)
		if int(letter) == '*' {
			imcfread_to_eol(fp)
			continue
		}
		if int(letter) != '#' {
			imcbug(libc.CString("%s"), "imc_load_color_table: # not found.")
			break
		}
		word = imcfread_word(fp)
		if C.strcasecmp(word, libc.CString("COLOR")) == 0 {
			for {
				if (func() *IMC_COLOR {
					color = new(IMC_COLOR)
					return color
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			imc_readcolor(color, fp)
			for {
				if first_imc_color == nil {
					first_imc_color = color
					last_imc_color = color
				} else {
					last_imc_color.Next = color
				}
				color.Next = nil
				if first_imc_color == color {
					color.Prev = nil
				} else {
					color.Prev = last_imc_color
				}
				last_imc_color = color
				if true {
					break
				}
			}
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_load_color_table: bad section: %s."), word)
			continue
		}
	}
	C.fclose(fp)
	fp = nil
}
func imc_savehelps() {
	var (
		fp   *C.FILE
		help *IMC_HELP_DATA
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "Couldn't write to IMC2 help file.")
		return
	}
	for help = first_imc_help; help != nil; help = help.Next {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#HELP\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Name %s\n", help.Name)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Perm %s\n", imcperm_names[help.Level])
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Text %s�\n", help.Text)
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#END\n")
	C.fclose(fp)
	fp = nil
}
func imc_readhelp(help *IMC_HELP_DATA, fp *C.FILE) {
	var (
		word      *byte
		hbuf      [4096]byte
		permvalue int
		fMatch    bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fp) != 0 {
				return "End"
			}
			return libc.GoString(imcfread_word(fp))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '*':
			fMatch = TRUE != 0
			imcfread_to_eol(fp)
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				return
			}
		case 'N':
			if C.strcasecmp(word, libc.CString("Name")) == 0 {
				help.Name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		case 'P':
			if C.strcasecmp(word, libc.CString("Perm")) == 0 {
				word = imcfread_word(fp)
				permvalue = get_imcpermvalue(word)
				if permvalue < 0 || permvalue > IMCPERM_IMP {
					imcbug(libc.CString("imc_readhelp: Command %s loaded with invalid permission. Set to Imp."), help.Name)
					help.Level = IMCPERM_IMP
				} else {
					help.Level = permvalue
				}
				fMatch = TRUE != 0
				break
			}
		case 'T':
			if C.strcasecmp(word, libc.CString("Text")) == 0 {
				var num int = 0
				for (func() byte {
					p := &hbuf[num]
					hbuf[num] = byte(int8(fgetc(fp)))
					return *p
				}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
					num++
				}
				hbuf[num] = '\x00'
				help.Text = C.strdup(&hbuf[0])
				fMatch = TRUE != 0
				break
			}
		}
		if !fMatch {
			imcbug(libc.CString("imc_readhelp: no match: %s"), word)
		}
	}
}
func imc_load_helps() {
	var (
		fp   *C.FILE
		help *IMC_HELP_DATA
	)
	first_imc_help = func() *IMC_HELP_DATA {
		last_imc_help = nil
		return last_imc_help
	}()
	imclog(libc.CString("%s"), "Loading IMC2 help file...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "No help file found.")
		return
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fp)
		if int(letter) == '*' {
			imcfread_to_eol(fp)
			continue
		}
		if int(letter) != '#' {
			imcbug(libc.CString("%s"), "imc_load_helps: # not found.")
			break
		}
		word = imcfread_word(fp)
		if C.strcasecmp(word, libc.CString("HELP")) == 0 {
			for {
				if (func() *IMC_HELP_DATA {
					help = new(IMC_HELP_DATA)
					return help
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			imc_readhelp(help, fp)
			for {
				if first_imc_help == nil {
					first_imc_help = help
					last_imc_help = help
				} else {
					last_imc_help.Next = help
				}
				help.Next = nil
				if first_imc_help == help {
					help.Prev = nil
				} else {
					help.Prev = last_imc_help
				}
				last_imc_help = help
				if true {
					break
				}
			}
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_load_helps: bad section: %s."), word)
			continue
		}
	}
	C.fclose(fp)
	fp = nil
}
func imc_savecommands() {
	var (
		fp    *C.FILE
		cmd   *IMC_CMD_DATA
		alias *IMC_ALIAS
	)
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "Couldn't write to IMC2 command file.")
		return
	}
	for cmd = first_imc_command; cmd != nil; cmd = cmd.Next {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#COMMAND\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Name      %s\n", cmd.Name)
		if cmd.Function != nil {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Code      %s\n", imc_funcname(cmd.Function))
		} else {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "Code      NULL\n")
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Perm      %s\n", imcperm_names[cmd.Level])
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Connected %d\n", cmd.Connected)
		for alias = cmd.First_alias; alias != nil; alias = alias.Next {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Alias     %s\n", alias.Name)
		}
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#END\n")
	C.fclose(fp)
	fp = nil
}
func imc_readcommand(cmd *IMC_CMD_DATA, fp *C.FILE) {
	var (
		alias     *IMC_ALIAS
		word      *byte
		permvalue int
		fMatch    bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fp) != 0 {
				return "End"
			}
			return libc.GoString(imcfread_word(fp))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '*':
			fMatch = TRUE != 0
			imcfread_to_eol(fp)
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				return
			}
		case 'A':
			if C.strcasecmp(word, libc.CString("Alias")) == 0 {
				for {
					if (func() *IMC_ALIAS {
						alias = new(IMC_ALIAS)
						return alias
					}()) == nil {
						imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
						abort()
					}
					if true {
						break
					}
				}
				alias.Name = imcfread_line(fp)
				for {
					if cmd.First_alias == nil {
						cmd.First_alias = alias
						cmd.Last_alias = alias
					} else {
						cmd.Last_alias.Next = alias
					}
					alias.Next = nil
					if cmd.First_alias == alias {
						alias.Prev = nil
					} else {
						alias.Prev = cmd.Last_alias
					}
					cmd.Last_alias = alias
					if true {
						break
					}
				}
				fMatch = TRUE != 0
				break
			}
		case 'C':
			if C.strcasecmp(word, libc.CString("Connected")) == 0 {
				cmd.Connected = imcfread_number(fp) != 0
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("Code")) == 0 {
				word = imcfread_word(fp)
				cmd.Function = imc_function(word)
				if cmd.Function == nil {
					imcbug(libc.CString("imc_readcommand: Command %s loaded with invalid function. Set to NULL."), cmd.Name)
				}
				fMatch = TRUE != 0
				break
			}
		case 'N':
			if C.strcasecmp(word, libc.CString("Name")) == 0 {
				cmd.Name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		case 'P':
			if C.strcasecmp(word, libc.CString("Perm")) == 0 {
				word = imcfread_word(fp)
				permvalue = get_imcpermvalue(word)
				if permvalue < 0 || permvalue > IMCPERM_IMP {
					imcbug(libc.CString("imc_readcommand: Command %s loaded with invalid permission. Set to Imp."), cmd.Name)
					cmd.Level = IMCPERM_IMP
				} else {
					cmd.Level = permvalue
				}
				fMatch = TRUE != 0
				break
			}
		}
		if !fMatch {
			imcbug(libc.CString("imc_readcommand: no match: %s"), word)
		}
	}
}
func imc_load_commands() bool {
	var (
		fp  *C.FILE
		cmd *IMC_CMD_DATA
	)
	first_imc_command = func() *IMC_CMD_DATA {
		last_imc_command = nil
		return last_imc_command
	}()
	imclog(libc.CString("%s"), "Loading IMC2 command table...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "No command table found.")
		return FALSE != 0
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fp)
		if int(letter) == '*' {
			imcfread_to_eol(fp)
			continue
		}
		if int(letter) != '#' {
			imcbug(libc.CString("%s"), "imc_load_commands: # not found.")
			break
		}
		word = imcfread_word(fp)
		if C.strcasecmp(word, libc.CString("COMMAND")) == 0 {
			for {
				if (func() *IMC_CMD_DATA {
					cmd = new(IMC_CMD_DATA)
					return cmd
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			imc_readcommand(cmd, fp)
			for {
				if first_imc_command == nil {
					first_imc_command = cmd
					last_imc_command = cmd
				} else {
					last_imc_command.Next = cmd
				}
				cmd.Next = nil
				if first_imc_command == cmd {
					cmd.Prev = nil
				} else {
					cmd.Prev = last_imc_command
				}
				last_imc_command = cmd
				if true {
					break
				}
			}
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_load_commands: bad section: %s."), word)
			continue
		}
	}
	C.fclose(fp)
	fp = nil
	return TRUE != 0
}
func imc_readucache(user *IMCUCACHE_DATA, fp *C.FILE) {
	var (
		word   *byte
		fMatch bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fp) != 0 {
				return "End"
			}
			return libc.GoString(imcfread_word(fp))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '*':
			fMatch = TRUE != 0
			imcfread_to_eol(fp)
		case 'N':
			if C.strcasecmp(word, libc.CString("Name")) == 0 {
				user.Name = imcfread_line(fp)
				fMatch = TRUE != 0
				break
			}
		case 'S':
			if C.strcasecmp(word, libc.CString("Sex")) == 0 {
				user.Gender = imcfread_number(fp)
				fMatch = TRUE != 0
				break
			}
		case 'T':
			if C.strcasecmp(word, libc.CString("Time")) == 0 {
				user.Time = int64(imcfread_number(fp))
				fMatch = TRUE != 0
				break
			}
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				return
			}
		}
		if !fMatch {
			imcbug(libc.CString("imc_readucache: no match: %s"), word)
		}
	}
}
func imc_load_ucache() {
	var (
		fp   *C.FILE
		user *IMCUCACHE_DATA
	)
	imclog(libc.CString("%s"), "Loading ucache data...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "No ucache data found.")
		return
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fp)
		if int(letter) == '*' {
			imcfread_to_eol(fp)
			continue
		}
		if int(letter) != '#' {
			imcbug(libc.CString("%s"), "imc_load_ucahe: # not found.")
			break
		}
		word = imcfread_word(fp)
		if C.strcasecmp(word, libc.CString("UCACHE")) == 0 {
			for {
				if (func() *IMCUCACHE_DATA {
					user = new(IMCUCACHE_DATA)
					return user
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			imc_readucache(user, fp)
			for {
				if first_imcucache == nil {
					first_imcucache = user
					last_imcucache = user
				} else {
					last_imcucache.Next = user
				}
				user.Next = nil
				if first_imcucache == user {
					user.Prev = nil
				} else {
					user.Prev = last_imcucache
				}
				last_imcucache = user
				if true {
					break
				}
			}
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_load_ucache: bad section: %s."), word)
			continue
		}
	}
	C.fclose(fp)
	fp = nil
	imc_prune_ucache()
	imcucache_clock = imc_time + 86400
}
func imc_save_config() {
	var fp *C.FILE
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s"), "Couldn't write to config file.")
		return
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "$IMCCONFIG\n\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "# %s config file.\n", this_imcmud.Versionid)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "# This file can now support the use of tildes in your strings.\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "# This information can be edited online using the 'imcconfig' command.\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "LocalName      %s\n", this_imcmud.Localname)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Autoconnect    %d\n", this_imcmud.Autoconnect)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "MinPlayerLevel %d\n", this_imcmud.Minlevel)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "MinImmLevel    %d\n", this_imcmud.Immlevel)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "AdminLevel     %d\n", this_imcmud.Adminlevel)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "Implevel       %d\n", this_imcmud.Implevel)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoName       %s\n", this_imcmud.Fullname)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoHost       %s\n", this_imcmud.Ihost)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoPort       %d\n", this_imcmud.Iport)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoEmail      %s\n", this_imcmud.Email)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoWWW        %s\n", this_imcmud.Www)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoBase       %s\n", this_imcmud.Base)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "InfoDetails    %s\n\n", this_imcmud.Details)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "# Your server connection information goes here.\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "# This information should be available from the network you plan to join.\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ServerAddr     %s\n", this_imcmud.Rhost)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ServerPort     %d\n", this_imcmud.Rport)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ClientPwd      %s\n", this_imcmud.Clientpw)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "ServerPwd      %s\n", this_imcmud.Serverpw)
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "#SHA256 auth: 0 = disabled, 1 = enabled\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "SHA256         %d\n", this_imcmud.Sha256)
	if this_imcmud.Sha256pass {
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "#Your server is expecting SHA256 authentication now. Do not remove this line unless told to do so.\n")
		stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "SHA256Pwd      %d\n", this_imcmud.Sha256pass)
	}
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "End\n\n")
	stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s", "$END\n")
	C.fclose(fp)
	fp = nil
}
func imcfread_config_file(fin *C.FILE) {
	var (
		word   *byte
		fMatch bool
	)
	for {
		word = libc.CString(func() string {
			if C.feof(fin) != 0 {
				return "end"
			}
			return libc.GoString(imcfread_word(fin))
		}())
		fMatch = FALSE != 0
		switch *word {
		case '#':
			fMatch = TRUE != 0
			imcfread_to_eol(fin)
		case 'A':
			if C.strcasecmp(word, libc.CString("Autoconnect")) == 0 {
				this_imcmud.Autoconnect = imcfread_number(fin) != 0
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("AdminLevel")) == 0 {
				this_imcmud.Adminlevel = imcfread_number(fin)
				fMatch = TRUE != 0
				break
			}
		case 'C':
			if C.strcasecmp(word, libc.CString("ClientPwd")) == 0 {
				this_imcmud.Clientpw = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
		case 'E':
			if C.strcasecmp(word, libc.CString("End")) == 0 {
				return
			}
		case 'I':
			if C.strcasecmp(word, libc.CString("Implevel")) == 0 {
				this_imcmud.Implevel = imcfread_number(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoName")) == 0 {
				this_imcmud.Fullname = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoHost")) == 0 {
				this_imcmud.Ihost = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoPort")) == 0 {
				this_imcmud.Iport = imcfread_number(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoEmail")) == 0 {
				this_imcmud.Email = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoWWW")) == 0 {
				this_imcmud.Www = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoBase")) == 0 {
				this_imcmud.Base = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("InfoDetails")) == 0 {
				this_imcmud.Details = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
		case 'L':
			if C.strcasecmp(word, libc.CString("LocalName")) == 0 {
				this_imcmud.Localname = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
		case 'M':
			if C.strcasecmp(word, libc.CString("MinImmLevel")) == 0 {
				this_imcmud.Immlevel = imcfread_number(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("MinPlayerLevel")) == 0 {
				this_imcmud.Minlevel = imcfread_number(fin)
				fMatch = TRUE != 0
				break
			}
		case 'R':
			if C.strcasecmp(word, libc.CString("RouterAddr")) == 0 {
				this_imcmud.Rhost = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("RouterPort")) == 0 {
				this_imcmud.Rport = uint16(int16(imcfread_number(fin)))
				fMatch = TRUE != 0
				break
			}
		case 'S':
			if C.strcasecmp(word, libc.CString("ServerPwd")) == 0 {
				this_imcmud.Serverpw = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ServerAddr")) == 0 {
				this_imcmud.Rhost = imcfread_line(fin)
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("ServerPort")) == 0 {
				this_imcmud.Rport = uint16(int16(imcfread_number(fin)))
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("SHA256")) == 0 {
				this_imcmud.Sha256 = imcfread_number(fin) != 0
				fMatch = TRUE != 0
				break
			}
			if C.strcasecmp(word, libc.CString("SHA256Pwd")) == 0 {
				this_imcmud.Sha256pass = imcfread_number(fin) != 0
				fMatch = TRUE != 0
				break
			}
		}
		if !fMatch {
			imcbug(libc.CString("%s: Bad keyword: %s\r\n"), libc.FuncName(), word)
		}
	}
}
func imc_read_config(desc int) bool {
	var (
		fin   *C.FILE
		cbase [1024]byte
	)
	if this_imcmud != nil {
		imc_delete_info()
	}
	this_imcmud = nil
	imclog(libc.CString("%s"), "Loading IMC2 network data...")
	if (func() *C.FILE {
		fin = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fin
	}()) == nil {
		imclog(libc.CString("%s"), "Can't open configuration file")
		imclog(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	for {
		var (
			letter int8
			word   *byte
		)
		letter = imcfread_letter(fin)
		if int(letter) == '#' {
			imcfread_to_eol(fin)
			continue
		}
		if int(letter) != '$' {
			imcbug(libc.CString("%s"), "imc_read_config: $ not found")
			break
		}
		word = imcfread_word(fin)
		if C.strcasecmp(word, libc.CString("IMCCONFIG")) == 0 && this_imcmud == nil {
			for {
				if (func() *SITEINFO {
					this_imcmud = new(SITEINFO)
					return this_imcmud
				}()) == nil {
					imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
					abort()
				}
				if true {
					break
				}
			}
			this_imcmud.Minlevel = 10
			this_imcmud.Immlevel = 101
			this_imcmud.Adminlevel = 113
			this_imcmud.Implevel = 115
			this_imcmud.Network = C.strdup(libc.CString("Unknown"))
			this_imcmud.Sha256 = TRUE != 0
			this_imcmud.Sha256pass = FALSE != 0
			this_imcmud.Desc = desc
			imcfread_config_file(fin)
			continue
		} else if C.strcasecmp(word, libc.CString("END")) == 0 {
			break
		} else {
			imcbug(libc.CString("imc_read_config: Bad section in config file: %s"), word)
			continue
		}
	}
	C.fclose(fin)
	fin = nil
	if this_imcmud == nil {
		imcbug(libc.CString("%s"), "imc_read_config: No server connection information!!")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Rhost == nil || this_imcmud.Clientpw == nil || this_imcmud.Serverpw == nil {
		imcbug(libc.CString("%s"), "imc_read_config: Missing required configuration info.")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Localname == nil || *this_imcmud.Localname == '\x00' {
		imcbug(libc.CString("%s"), "imc_read_config: Mud name not loaded in configuration file.")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Fullname == nil || *this_imcmud.Fullname == '\x00' {
		imcbug(libc.CString("%s"), "imc_read_config: Missing InfoName parameter in configuration file.")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Ihost == nil || *this_imcmud.Ihost == '\x00' {
		imcbug(libc.CString("%s"), "imc_read_config: Missing InfoHost parameter in configuration file.")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Email == nil || *this_imcmud.Email == '\x00' {
		imcbug(libc.CString("%s"), "imc_read_config: Missing InfoEmail parameter in configuration file.")
		imcbug(libc.CString("%s"), "Network configuration aborted.")
		return FALSE != 0
	}
	if this_imcmud.Base == nil || *this_imcmud.Base == '\x00' {
		this_imcmud.Base = C.strdup(libc.CString("Unknown Codebase"))
	}
	if this_imcmud.Www == nil || *this_imcmud.Www == '\x00' {
		this_imcmud.Www = C.strdup(libc.CString("Not specified"))
	}
	if this_imcmud.Details == nil || *this_imcmud.Details == '\x00' {
		this_imcmud.Details = C.strdup(libc.CString("No details provided."))
	}
	if this_imcmud.Versionid == nil {
		stdio.Snprintf(&cbase[0], SMST, "%s%s", IMC_VERSION_STRING, this_imcmud.Base)
		this_imcmud.Versionid = C.strdup(&cbase[0])
	}
	return TRUE != 0
}
func parse_who_header(head *byte) *byte {
	var (
		newhead [4096]byte
		iport   [1024]byte
	)
	stdio.Snprintf(&iport[0], SMST, "%d", this_imcmud.Iport)
	strlcpy(&newhead[0], head, LGST)
	strlcpy(&newhead[0], imcstrrep(&newhead[0], libc.CString("<%mudfullname%>"), this_imcmud.Fullname), LGST)
	strlcpy(&newhead[0], imcstrrep(&newhead[0], libc.CString("<%mudtelnet%>"), this_imcmud.Ihost), LGST)
	strlcpy(&newhead[0], imcstrrep(&newhead[0], libc.CString("<%mudport%>"), &iport[0]), LGST)
	strlcpy(&newhead[0], imcstrrep(&newhead[0], libc.CString("<%mudurl%>"), this_imcmud.Www), LGST)
	return &newhead[0]
}
func parse_who_tail(tail *byte) *byte {
	var (
		newtail [4096]byte
		iport   [1024]byte
	)
	stdio.Snprintf(&iport[0], SMST, "%d", this_imcmud.Iport)
	strlcpy(&newtail[0], tail, LGST)
	strlcpy(&newtail[0], imcstrrep(&newtail[0], libc.CString("<%mudfullname%>"), this_imcmud.Fullname), LGST)
	strlcpy(&newtail[0], imcstrrep(&newtail[0], libc.CString("<%mudtelnet%>"), this_imcmud.Ihost), LGST)
	strlcpy(&newtail[0], imcstrrep(&newtail[0], libc.CString("<%mudport%>"), &iport[0]), LGST)
	strlcpy(&newtail[0], imcstrrep(&newtail[0], libc.CString("<%mudurl%>"), this_imcmud.Www), LGST)
	return &newtail[0]
}
func imc_delete_who_template() {
	for {
		if whot.Head != nil {
			libc.Free(unsafe.Pointer(whot.Head))
			whot.Head = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Plrheader != nil {
			libc.Free(unsafe.Pointer(whot.Plrheader))
			whot.Plrheader = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Immheader != nil {
			libc.Free(unsafe.Pointer(whot.Immheader))
			whot.Immheader = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Plrline != nil {
			libc.Free(unsafe.Pointer(whot.Plrline))
			whot.Plrline = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Immline != nil {
			libc.Free(unsafe.Pointer(whot.Immline))
			whot.Immline = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Tail != nil {
			libc.Free(unsafe.Pointer(whot.Tail))
			whot.Tail = nil
		}
		if true {
			break
		}
	}
	for {
		if whot.Master != nil {
			libc.Free(unsafe.Pointer(whot.Master))
			whot.Master = nil
		}
		if true {
			break
		}
	}
	for {
		if whot != nil {
			libc.Free(unsafe.Pointer(whot))
			whot = nil
		}
		if true {
			break
		}
	}
}
func imc_load_who_template() {
	var (
		fp   *C.FILE
		hbuf [4096]byte
		word *byte
		num  int
	)
	imclog(libc.CString("%s"), "Loading IMC2 who template...")
	if (func() *C.FILE {
		fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
		return fp
	}()) == nil {
		imclog(libc.CString("%s: Unable to load template file for imcwho"), libc.FuncName())
		whot = nil
		return
	}
	if whot != nil {
		imc_delete_who_template()
	}
	for {
		if (func() *WHO_TEMPLATE {
			whot = new(WHO_TEMPLATE)
			return whot
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	for {
		word = imcfread_word(fp)
		hbuf[0] = '\x00'
		num = 0
		if C.strcasecmp(word, libc.CString("Head:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Head = C.strdup(parse_who_header(&hbuf[0]))
		} else if C.strcasecmp(word, libc.CString("Tail:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Tail = C.strdup(parse_who_tail(&hbuf[0]))
		} else if C.strcasecmp(word, libc.CString("Plrline:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Plrline = C.strdup(&hbuf[0])
		} else if C.strcasecmp(word, libc.CString("Immline:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Immline = C.strdup(&hbuf[0])
		} else if C.strcasecmp(word, libc.CString("Immheader:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Immheader = C.strdup(&hbuf[0])
		} else if C.strcasecmp(word, libc.CString("Plrheader:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Plrheader = C.strdup(&hbuf[0])
		} else if C.strcasecmp(word, libc.CString("Master:")) == 0 {
			for (func() byte {
				p := &hbuf[num]
				hbuf[num] = byte(int8(fgetc(fp)))
				return *p
			}()) != math.MaxUint8 && hbuf[num] != 'z' && num < (int(LGST-2)) {
				num++
			}
			hbuf[num] = '\x00'
			whot.Master = C.strdup(&hbuf[0])
		}
		if C.feof(fp) != 0 {
			break
		}
	}
	C.fclose(fp)
	fp = nil
}
func imc_load_templates() {
	imc_load_who_template()
}
func ipv4_connect() int {
	var (
		sa    sockaddr_in
		hostp *hostent
		r     int
		desc  int = -1
	)
	sa = sockaddr_in{}
	sa.Sin_family = PF_INET
	if inet_aton(this_imcmud.Rhost, &sa.Sin_addr) == 0 {
		hostp = gethostbyname(this_imcmud.Rhost)
		if hostp == nil {
			imclog(libc.CString("%s"), "imc_connect_to: Cannot resolve server hostname.")
			imc_shutdown(FALSE != 0)
			return -1
		}
		libc.MemCpy(unsafe.Pointer(&sa.Sin_addr), unsafe.Pointer(*(**byte)(unsafe.Add(unsafe.Pointer(hostp.H_addr_list), unsafe.Sizeof((*byte)(nil))*0))), hostp.H_length)
	}
	sa.Sin_port = in_port_t(htons(uint16(this_imcmud.Rport)))
	desc = socket(PF_INET, SOCK_STREAM, 0)
	if desc < 0 {
		C.perror(libc.CString("socket"))
		return -1
	}
	r = fcntl(desc, F_GETFL, 0)
	if r < 0 || fcntl(desc, F_SETFL, O_NONBLOCK|r) < 0 {
		C.perror(libc.CString("imc_connect: fcntl"))
		close_(desc)
		return -1
	}
	if connect(desc, (*sockaddr)(unsafe.Pointer(&sa)), int(unsafe.Sizeof(sockaddr_in{}))) == -1 {
		if (*__errno_location()) != EINPROGRESS {
			imclog(libc.CString("%s: Failed connect: Error %d: %s"), libc.FuncName(), *__errno_location(), C.strerror(*__errno_location()))
			C.perror(libc.CString("connect"))
			close_(desc)
			return -1
		}
	}
	return desc
}
func imc_server_connect() bool {
	var (
		buf  [4096]byte
		desc int = 0
	)
	if this_imcmud == nil {
		imcbug(libc.CString("%s"), "No connection data loaded")
		return FALSE != 0
	}
	if int(this_imcmud.State) != IMC_AUTH1 {
		imcbug(libc.CString("%s"), "Connection is not in proper state.")
		return FALSE != 0
	}
	if this_imcmud.Desc > 0 {
		imcbug(libc.CString("%s"), "Already connected")
		return FALSE != 0
	}
	desc = ipv4_connect()
	if desc < 1 {
		return FALSE != 0
	}
	imclog(libc.CString("%s"), "Connecting to server.")
	this_imcmud.State = uint16(int16(IMC_AUTH2))
	this_imcmud.Desc = desc
	this_imcmud.Inbuf[0] = '\x00'
	this_imcmud.Outsize = 1000
	for {
		if (func() *byte {
			p := &this_imcmud.Outbuf
			this_imcmud.Outbuf = (*byte)(unsafe.Pointer(&make([]int8, int(this_imcmud.Outsize))[0]))
			return *p
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	if this_imcmud.Sha256 {
		if !this_imcmud.Sha256pass {
			stdio.Snprintf(&buf[0], LGST, "PW %s %s version=%d autosetup %s SHA256", this_imcmud.Localname, this_imcmud.Clientpw, IMC_VERSION, this_imcmud.Serverpw)
		} else {
			stdio.Snprintf(&buf[0], LGST, "SHA256-AUTH-REQ %s", this_imcmud.Localname)
		}
	} else {
		stdio.Snprintf(&buf[0], LGST, "PW %s %s version=%d autosetup %s", this_imcmud.Localname, this_imcmud.Clientpw, IMC_VERSION, this_imcmud.Serverpw)
	}
	imc_write_buffer(&buf[0])
	return TRUE != 0
}
func imc_delete_templates() {
	imc_delete_who_template()
}
func free_imcdata(complete bool) {
	var (
		p           *REMOTEINFO
		pnext       *REMOTEINFO
		ban         *IMC_BAN
		ban_next    *IMC_BAN
		ucache      *IMCUCACHE_DATA
		next_ucache *IMCUCACHE_DATA
		cmd         *IMC_CMD_DATA
		cmd_next    *IMC_CMD_DATA
		alias       *IMC_ALIAS
		alias_next  *IMC_ALIAS
		help        *IMC_HELP_DATA
		help_next   *IMC_HELP_DATA
		color       *IMC_COLOR
		color_next  *IMC_COLOR
		ph          *IMC_PHANDLER
		ph_next     *IMC_PHANDLER
		c           *IMC_CHANNEL
		c_next      *IMC_CHANNEL
	)
	for c = first_imc_channel; c != nil; c = c_next {
		c_next = c.Next
		imc_freechan(c)
	}
	for p = first_rinfo; p != nil; p = pnext {
		pnext = p.Next
		imc_delete_reminfo(p)
	}
	for ban = first_imc_ban; ban != nil; ban = ban_next {
		ban_next = ban.Next
		imc_freeban(ban)
	}
	for ucache = first_imcucache; ucache != nil; ucache = next_ucache {
		next_ucache = ucache.Next
		for {
			if ucache.Name != nil {
				libc.Free(unsafe.Pointer(ucache.Name))
				ucache.Name = nil
			}
			if true {
				break
			}
		}
		for {
			if ucache.Prev == nil {
				first_imcucache = ucache.Next
				if first_imcucache != nil {
					first_imcucache.Prev = nil
				}
			} else {
				ucache.Prev.Next = ucache.Next
			}
			if ucache.Next == nil {
				last_imcucache = ucache.Prev
				if last_imcucache != nil {
					last_imcucache.Next = nil
				}
			} else {
				ucache.Next.Prev = ucache.Prev
			}
			if true {
				break
			}
		}
		for {
			if ucache != nil {
				libc.Free(unsafe.Pointer(ucache))
				ucache = nil
			}
			if true {
				break
			}
		}
	}
	if complete {
		imc_delete_templates()
		for cmd = first_imc_command; cmd != nil; cmd = cmd_next {
			cmd_next = cmd.Next
			for alias = cmd.First_alias; alias != nil; alias = alias_next {
				alias_next = alias.Next
				for {
					if alias.Name != nil {
						libc.Free(unsafe.Pointer(alias.Name))
						alias.Name = nil
					}
					if true {
						break
					}
				}
				for {
					if alias.Prev == nil {
						cmd.First_alias = alias.Next
						if cmd.First_alias != nil {
							cmd.First_alias.Prev = nil
						}
					} else {
						alias.Prev.Next = alias.Next
					}
					if alias.Next == nil {
						cmd.Last_alias = alias.Prev
						if cmd.Last_alias != nil {
							cmd.Last_alias.Next = nil
						}
					} else {
						alias.Next.Prev = alias.Prev
					}
					if true {
						break
					}
				}
				for {
					if alias != nil {
						libc.Free(unsafe.Pointer(alias))
						alias = nil
					}
					if true {
						break
					}
				}
			}
			for {
				if cmd.Name != nil {
					libc.Free(unsafe.Pointer(cmd.Name))
					cmd.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if cmd.Prev == nil {
					first_imc_command = cmd.Next
					if first_imc_command != nil {
						first_imc_command.Prev = nil
					}
				} else {
					cmd.Prev.Next = cmd.Next
				}
				if cmd.Next == nil {
					last_imc_command = cmd.Prev
					if last_imc_command != nil {
						last_imc_command.Next = nil
					}
				} else {
					cmd.Next.Prev = cmd.Prev
				}
				if true {
					break
				}
			}
			for {
				if cmd != nil {
					libc.Free(unsafe.Pointer(cmd))
					cmd = nil
				}
				if true {
					break
				}
			}
		}
		for help = first_imc_help; help != nil; help = help_next {
			help_next = help.Next
			for {
				if help.Name != nil {
					libc.Free(unsafe.Pointer(help.Name))
					help.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if help.Text != nil {
					libc.Free(unsafe.Pointer(help.Text))
					help.Text = nil
				}
				if true {
					break
				}
			}
			for {
				if help.Prev == nil {
					first_imc_help = help.Next
					if first_imc_help != nil {
						first_imc_help.Prev = nil
					}
				} else {
					help.Prev.Next = help.Next
				}
				if help.Next == nil {
					last_imc_help = help.Prev
					if last_imc_help != nil {
						last_imc_help.Next = nil
					}
				} else {
					help.Next.Prev = help.Prev
				}
				if true {
					break
				}
			}
			for {
				if help != nil {
					libc.Free(unsafe.Pointer(help))
					help = nil
				}
				if true {
					break
				}
			}
		}
		for color = first_imc_color; color != nil; color = color_next {
			color_next = color.Next
			for {
				if color.Name != nil {
					libc.Free(unsafe.Pointer(color.Name))
					color.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if color.Mudtag != nil {
					libc.Free(unsafe.Pointer(color.Mudtag))
					color.Mudtag = nil
				}
				if true {
					break
				}
			}
			for {
				if color.Imctag != nil {
					libc.Free(unsafe.Pointer(color.Imctag))
					color.Imctag = nil
				}
				if true {
					break
				}
			}
			for {
				if color.Prev == nil {
					first_imc_color = color.Next
					if first_imc_color != nil {
						first_imc_color.Prev = nil
					}
				} else {
					color.Prev.Next = color.Next
				}
				if color.Next == nil {
					last_imc_color = color.Prev
					if last_imc_color != nil {
						last_imc_color.Next = nil
					}
				} else {
					color.Next.Prev = color.Prev
				}
				if true {
					break
				}
			}
			for {
				if color != nil {
					libc.Free(unsafe.Pointer(color))
					color = nil
				}
				if true {
					break
				}
			}
		}
		for ph = first_phandler; ph != nil; ph = ph_next {
			ph_next = ph.Next
			for {
				if ph.Name != nil {
					libc.Free(unsafe.Pointer(ph.Name))
					ph.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if ph.Prev == nil {
					first_phandler = ph.Next
					if first_phandler != nil {
						first_phandler.Prev = nil
					}
				} else {
					ph.Prev.Next = ph.Next
				}
				if ph.Next == nil {
					last_phandler = ph.Prev
					if last_phandler != nil {
						last_phandler.Next = nil
					}
				} else {
					ph.Next.Prev = ph.Prev
				}
				if true {
					break
				}
			}
			for {
				if ph != nil {
					libc.Free(unsafe.Pointer(ph))
					ph = nil
				}
				if true {
					break
				}
			}
		}
	}
}
func imc_hotboot() {
	var fp *C.FILE
	if this_imcmud != nil && int(this_imcmud.State) == IMC_ONLINE {
		if (func() *C.FILE {
			fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "w")))
			return fp
		}()) == nil {
			imcbug(libc.CString("%s: Unable to open IMC hotboot file for write."), libc.FuncName())
		} else {
			stdio.Fprintf((*stdio.File)(unsafe.Pointer(fp)), "%s %s\n", func() *byte {
				if this_imcmud.Network != nil {
					return this_imcmud.Network
				}
				return libc.CString("Unknown")
			}(), func() *byte {
				if this_imcmud.Servername != nil {
					return this_imcmud.Servername
				}
				return libc.CString("Unknown")
			}())
			C.fclose(fp)
			fp = nil
			imc_savehistory()
		}
	}
}
func imc_shutdown(reconnect bool) {
	if this_imcmud != nil && int(this_imcmud.State) == IMC_OFFLINE {
		return
	}
	imclog(libc.CString("%s"), "Shutting down network.")
	if this_imcmud.Desc > 0 {
		close_(this_imcmud.Desc)
	}
	this_imcmud.Desc = -1
	imc_savehistory()
	free_imcdata(FALSE != 0)
	this_imcmud.State = uint16(int16(IMC_OFFLINE))
	if reconnect {
		imcwait = 100
		imclog(libc.CString("%s"), "Connection to server was lost. Reconnecting in approximately 20 seconds.")
	}
}
func imc_startup_network(connected bool) bool {
	imclog(libc.CString("%s"), "IMC2 Network Initializing...")
	if connected {
		var (
			fp      *C.FILE
			netname [1024]byte
			server  [1024]byte
		)
		if (func() *C.FILE {
			fp = (*C.FILE)(unsafe.Pointer(stdio.FOpen(IMC_DIR, "r")))
			return fp
		}()) == nil {
			imcbug(libc.CString("%s: Unable to load IMC hotboot file."), libc.FuncName())
		} else {
			unlink(libc.CString(IMC_DIR))
			__isoc99_fscanf(fp, libc.CString("%s %s\n"), &netname[0], &server[0])
			for {
				if this_imcmud.Network != nil {
					libc.Free(unsafe.Pointer(this_imcmud.Network))
					this_imcmud.Network = nil
				}
				if true {
					break
				}
			}
			this_imcmud.Network = C.strdup(&netname[0])
			for {
				if this_imcmud.Servername != nil {
					libc.Free(unsafe.Pointer(this_imcmud.Servername))
					this_imcmud.Servername = nil
				}
				if true {
					break
				}
			}
			this_imcmud.Servername = C.strdup(&server[0])
			C.fclose(fp)
			fp = nil
		}
		this_imcmud.State = uint16(int16(IMC_ONLINE))
		this_imcmud.Inbuf[0] = '\x00'
		this_imcmud.Outsize = IMC_BUFF_SIZE
		for {
			if (func() *byte {
				p := &this_imcmud.Outbuf
				this_imcmud.Outbuf = (*byte)(unsafe.Pointer(&make([]int8, int(this_imcmud.Outsize))[0]))
				return *p
			}()) == nil {
				imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
				abort()
			}
			if true {
				break
			}
		}
		imc_request_keepalive()
		imc_firstrefresh()
		return TRUE != 0
	}
	this_imcmud.State = uint16(int16(IMC_AUTH1))
	if !imc_server_connect() {
		this_imcmud.State = uint16(int16(IMC_OFFLINE))
		return FALSE != 0
	}
	return TRUE != 0
}
func imc_startup(force bool, desc int, connected bool) {
	imcwait = 0
	if this_imcmud != nil && int(this_imcmud.State) > IMC_OFFLINE {
		imclog(libc.CString("%s: Network startup called when already engaged!"), libc.FuncName())
		return
	}
	imc_time = C.time(nil)
	imc_sequencenumber = uint(imc_time)
	if first_imc_command == nil {
		if !imc_load_commands() {
			imcbug(libc.CString("%s: Unable to load command table!"), libc.FuncName())
			return
		}
	}
	if !imc_read_config(desc) {
		return
	}
	imc_register_default_packets()
	if first_imc_help == nil {
		imc_load_helps()
	}
	if first_imc_color == nil {
		imc_load_color_table()
	}
	if whot == nil {
		imc_load_templates()
	}
	if !this_imcmud.Autoconnect && !force && !connected || connected && this_imcmud.Desc < 1 {
		imclog(libc.CString("%s"), "IMC2 network data loaded. Autoconnect not set. IMC2 will need to be connected manually.")
		return
	} else {
		imclog(libc.CString("%s"), "IMC2 network data loaded.")
	}
	if this_imcmud.Autoconnect || force || connected {
		if imc_startup_network(connected) {
			imc_loadchannels()
			imc_loadhistory()
			imc_readbans()
			imc_load_ucache()
		}
	}
}
func imccommand(ch *char_data, argument *byte) {
	var (
		cmd   [1024]byte
		chan_ [1024]byte
		to    [1024]byte
		p     *IMC_PACKET
		c     *IMC_CHANNEL
	)
	argument = imcone_argument(argument, &cmd[0])
	argument = imcone_argument(argument, &chan_[0])
	if cmd[0] == 0 || chan_[0] == 0 {
		imc_to_char(libc.CString("Syntax: imccommand <command> <server:channel> [<data..>]\r\n"), ch)
		imc_to_char(libc.CString("Command access will depend on your privledges and what each individual server allows.\r\n"), ch)
		return
	}
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil && C.strcasecmp(&cmd[0], libc.CString("create")) != 0 {
		imc_printf(ch, libc.CString("There is no channel called %s known.\r\n"), &chan_[0])
		return
	}
	stdio.Snprintf(&to[0], SMST, "IMC@%s", func() *byte {
		if c != nil {
			return imc_channel_mudof(c.Name)
		}
		return imc_channel_mudof(&chan_[0])
	}())
	p = imc_newpacket(GET_NAME(ch), libc.CString("ice-cmd"), &to[0])
	imc_addtopacket(p, libc.CString("channel=%s"), func() *byte {
		if c != nil {
			return c.Name
		}
		return &chan_[0]
	}())
	imc_addtopacket(p, libc.CString("command=%s"), &cmd[0])
	if argument != nil && *argument != '\x00' {
		imc_addtopacket(p, libc.CString("data=%s"), argument)
	}
	imc_write_packet(p)
	imc_to_char(libc.CString("Command sent.\r\n"), ch)
}
func verify_format(fmt *byte, sneed int16) bool {
	var (
		c *byte
		i int = 0
	)
	c = fmt
	for (func() *byte {
		c = C.strchr(c, '%')
		return c
	}()) != nil {
		if *((*byte)(unsafe.Add(unsafe.Pointer(c), 1))) == '%' {
			c = (*byte)(unsafe.Add(unsafe.Pointer(c), 2))
			continue
		}
		if *((*byte)(unsafe.Add(unsafe.Pointer(c), 1))) != 's' {
			return FALSE != 0
		}
		c = (*byte)(unsafe.Add(unsafe.Pointer(c), 1))
		i++
	}
	if i != int(sneed) {
		return FALSE != 0
	}
	return TRUE != 0
}
func imcsetup(ch *char_data, argument *byte) {
	var (
		imccmd [1024]byte
		chan_  [1024]byte
		arg1   [1024]byte
		buf    [4096]byte
		c      *IMC_CHANNEL = nil
		x      int
		all    bool = FALSE != 0
	)
	argument = imcone_argument(argument, &imccmd[0])
	argument = imcone_argument(argument, &chan_[0])
	argument = imcone_argument(argument, &arg1[0])
	if imccmd[0] == '\x00' || chan_[0] == '\x00' {
		imc_to_char(libc.CString("Syntax: imcsetup <command> <channel> [<data..>]\r\n"), ch)
		imc_to_char(libc.CString("Where 'command' is one of the following:\r\n"), ch)
		imc_to_char(libc.CString("delete rename perm regformat emoteformat socformat\r\n\r\n"), ch)
		imc_to_char(libc.CString("Where 'channel' is one of the following:\r\n"), ch)
		for c = first_imc_channel; c != nil; c = c.Next {
			if c.Local_name != nil && *c.Local_name != '\x00' {
				imc_printf(ch, libc.CString("%s "), c.Local_name)
			} else {
				imc_printf(ch, libc.CString("%s "), c.Name)
			}
		}
		imc_to_char(libc.CString("\r\n"), ch)
		return
	}
	if C.strcasecmp(&chan_[0], libc.CString("all")) == 0 {
		all = TRUE != 0
	} else {
		if (func() *IMC_CHANNEL {
			c = imc_findchannel(&chan_[0])
			return c
		}()) == nil {
			imc_to_char(libc.CString("Unknown channel.\r\n"), ch)
			return
		}
	}
	if c != nil && int(c.Level) > ch.Player_specials.Imcchardata.Imcperm {
		imc_to_char(libc.CString("You cannot modify that channel."), ch)
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("delete")) == 0 {
		if all {
			imc_to_char(libc.CString("You cannot perform a delete all on channels.\r\n"), ch)
			return
		}
		for {
			if c.Local_name != nil {
				libc.Free(unsafe.Pointer(c.Local_name))
				c.Local_name = nil
			}
			if true {
				break
			}
		}
		for {
			if c.Regformat != nil {
				libc.Free(unsafe.Pointer(c.Regformat))
				c.Regformat = nil
			}
			if true {
				break
			}
		}
		for {
			if c.Emoteformat != nil {
				libc.Free(unsafe.Pointer(c.Emoteformat))
				c.Emoteformat = nil
			}
			if true {
				break
			}
		}
		for {
			if c.Socformat != nil {
				libc.Free(unsafe.Pointer(c.Socformat))
				c.Socformat = nil
			}
			if true {
				break
			}
		}
		for x = 0; x < MAX_IMCHISTORY; x++ {
			for {
				if (c.History[x]) != nil {
					libc.Free(unsafe.Pointer(c.History[x]))
					c.History[x] = nil
				}
				if true {
					break
				}
			}
		}
		imc_to_char(libc.CString("Channel is no longer locally configured.\r\n"), ch)
		if !c.Refreshed {
			imc_freechan(c)
		}
		imc_save_channels()
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("rename")) == 0 {
		if all {
			imc_to_char(libc.CString("You cannot perform a rename all on channels.\r\n"), ch)
			return
		}
		if arg1[0] == '\x00' {
			imc_to_char(libc.CString("Missing 'newname' argument for 'imcsetup rename'\r\n"), ch)
			imc_to_char(libc.CString("Syntax: imcsetup rename <local channel> <newname>\r\n"), ch)
			return
		}
		if imc_findchannel(&arg1[0]) != nil {
			imc_to_char(libc.CString("New channel name already exists.\r\n"), ch)
			return
		}
		stdio.Snprintf(&buf[0], LGST, "Renamed channel '%s' to '%s'.\r\n", c.Local_name, &arg1[0])
		for {
			if c.Local_name != nil {
				libc.Free(unsafe.Pointer(c.Local_name))
				c.Local_name = nil
			}
			if true {
				break
			}
		}
		c.Local_name = C.strdup(&arg1[0])
		imc_to_char(&buf[0], ch)
		imcformat_channel(ch, c, 4, FALSE != 0)
		imc_save_channels()
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("resetformats")) == 0 {
		if all {
			imcformat_channel(ch, nil, 4, TRUE != 0)
			imc_to_char(libc.CString("All channel formats have been reset to default.\r\n"), ch)
		} else {
			imcformat_channel(ch, c, 4, FALSE != 0)
			imc_to_char(libc.CString("The formats for this channel have been reset to default.\r\n"), ch)
		}
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("regformat")) == 0 {
		if arg1[0] == 0 {
			imc_to_char(libc.CString("Syntax: imcsetup regformat <localchannel|all> <string>\r\n"), ch)
			return
		}
		if !verify_format(&arg1[0], 2) {
			imc_to_char(libc.CString("Bad format - must contain exactly 2 %s's.\r\n"), ch)
			return
		}
		if all {
			for c = first_imc_channel; c != nil; c = c.Next {
				for {
					if c.Regformat != nil {
						libc.Free(unsafe.Pointer(c.Regformat))
						c.Regformat = nil
					}
					if true {
						break
					}
				}
				c.Regformat = C.strdup(&arg1[0])
			}
			imc_to_char(libc.CString("All channel regular formats have been updated.\r\n"), ch)
		} else {
			for {
				if c.Regformat != nil {
					libc.Free(unsafe.Pointer(c.Regformat))
					c.Regformat = nil
				}
				if true {
					break
				}
			}
			c.Regformat = C.strdup(&arg1[0])
			imc_to_char(libc.CString("The regular format for this channel has been changed successfully.\r\n"), ch)
		}
		imc_save_channels()
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("emoteformat")) == 0 {
		if arg1[0] == 0 {
			imc_to_char(libc.CString("Syntax: imcsetup emoteformat <localchannel|all> <string>\r\n"), ch)
			return
		}
		if !verify_format(&arg1[0], 2) {
			imc_to_char(libc.CString("Bad format - must contain exactly 2 %s's.\r\n"), ch)
			return
		}
		if all {
			for c = first_imc_channel; c != nil; c = c.Next {
				for {
					if c.Emoteformat != nil {
						libc.Free(unsafe.Pointer(c.Emoteformat))
						c.Emoteformat = nil
					}
					if true {
						break
					}
				}
				c.Emoteformat = C.strdup(&arg1[0])
			}
			imc_to_char(libc.CString("All channel emote formats have been updated.\r\n"), ch)
		} else {
			for {
				if c.Emoteformat != nil {
					libc.Free(unsafe.Pointer(c.Emoteformat))
					c.Emoteformat = nil
				}
				if true {
					break
				}
			}
			c.Emoteformat = C.strdup(&arg1[0])
			imc_to_char(libc.CString("The emote format for this channel has been changed successfully.\r\n"), ch)
		}
		imc_save_channels()
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("socformat")) == 0 {
		if arg1[0] == 0 {
			imc_to_char(libc.CString("Syntax: imcsetup socformat <localchannel|all> <string>\r\n"), ch)
			return
		}
		if !verify_format(&arg1[0], 1) {
			imc_to_char(libc.CString("Bad format - must contain exactly 1 %s.\r\n"), ch)
			return
		}
		if all {
			for c = first_imc_channel; c != nil; c = c.Next {
				for {
					if c.Socformat != nil {
						libc.Free(unsafe.Pointer(c.Socformat))
						c.Socformat = nil
					}
					if true {
						break
					}
				}
				c.Socformat = C.strdup(&arg1[0])
			}
			imc_to_char(libc.CString("All channel social formats have been updated.\r\n"), ch)
		} else {
			for {
				if c.Socformat != nil {
					libc.Free(unsafe.Pointer(c.Socformat))
					c.Socformat = nil
				}
				if true {
					break
				}
			}
			c.Socformat = C.strdup(&arg1[0])
			imc_to_char(libc.CString("The social format for this channel has been changed successfully.\r\n"), ch)
		}
		imc_save_channels()
		return
	}
	if C.strcasecmp(&imccmd[0], libc.CString("perm")) == 0 || C.strcasecmp(&imccmd[0], libc.CString("permission")) == 0 || C.strcasecmp(&imccmd[0], libc.CString("level")) == 0 {
		var permvalue int = -1
		if all {
			imc_to_char(libc.CString("You cannot do a permissions all for channels.\r\n"), ch)
			return
		}
		if arg1[0] == 0 {
			imc_to_char(libc.CString("Syntax: imcsetup perm <localchannel> <permission>\r\n"), ch)
			return
		}
		permvalue = get_imcpermvalue(&arg1[0])
		if permvalue < 0 || permvalue > IMCPERM_IMP {
			imc_to_char(libc.CString("Unacceptable permission setting.\r\n"), ch)
			return
		}
		if permvalue > ch.Player_specials.Imcchardata.Imcperm {
			imc_to_char(libc.CString("You cannot set a permission higher than your own.\r\n"), ch)
			return
		}
		c.Level = int16(permvalue)
		imc_to_char(libc.CString("Channel permissions changed.\r\n"), ch)
		imc_save_channels()
		return
	}
	imcsetup(ch, libc.CString(""))
}
func imcchanlist(ch *char_data, argument *byte) {
	var (
		c     *IMC_CHANNEL = nil
		count int          = 0
		col   int8         = 'C'
	)
	if first_imc_channel == nil {
		imc_to_char(libc.CString("~WThere are no known channels on this network.\r\n"), ch)
		return
	}
	if argument != nil && *argument != '\x00' {
		if (func() *IMC_CHANNEL {
			c = imc_findchannel(argument)
			return c
		}()) == nil {
			imc_printf(ch, libc.CString("There is no channel called %s here.\r\n"), argument)
			return
		}
	}
	if c != nil {
		imc_printf(ch, libc.CString("~WChannel  : %s\r\n\r\n"), c.Name)
		imc_printf(ch, libc.CString("~cLocalname: ~w%s\r\n"), c.Local_name)
		imc_printf(ch, libc.CString("~cPerms    : ~w%s\r\n"), imcperm_names[c.Level])
		imc_printf(ch, libc.CString("~cPolicy   : %s\r\n"), func() string {
			if c.Open {
				return "~gOpen"
			}
			return "~yPrivate"
		}())
		imc_printf(ch, libc.CString("~cRegFormat: ~w%s\r\n"), c.Regformat)
		imc_printf(ch, libc.CString("~cEmoFormat: ~w%s\r\n"), c.Emoteformat)
		imc_printf(ch, libc.CString("~cSocFormat: ~w%s\r\n\r\n"), c.Socformat)
		imc_printf(ch, libc.CString("~cOwner    : ~w%s\r\n"), c.Owner)
		imc_printf(ch, libc.CString("~cOperators: ~w%s\r\n"), c.Operators)
		imc_printf(ch, libc.CString("~cInvite   : ~w%s\r\n"), c.Invited)
		imc_printf(ch, libc.CString("~cExclude  : ~w%s\r\n"), c.Excluded)
		return
	}
	imc_printf(ch, libc.CString("~c%-15s ~C%-15s ~B%-15s ~b%-7s ~!%s\r\n\r\n"), "Name", "Local name", "Owner", "Perm", "Policy")
	for c = first_imc_channel; c != nil; c = c.Next {
		if ch.Player_specials.Imcchardata.Imcperm < int(c.Level) {
			continue
		}
		if c.Local_name != nil {
			if !imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
				col = 'R'
			} else {
				col = 'C'
			}
		}
		imc_printf(ch, libc.CString("~c%-15.15s ~%c%-*.*s ~B%-15.15s ~b%-7s %s\r\n"), c.Name, col, func() int {
			if c.Local_name != nil {
				return 15
			}
			return 17
		}(), func() int {
			if c.Local_name != nil {
				return 15
			}
			return 17
		}(), func() *byte {
			if c.Local_name != nil {
				return c.Local_name
			}
			return libc.CString("~Y(not local)  ")
		}(), c.Owner, imcperm_names[c.Level], func() string {
			if c.Refreshed {
				if c.Open {
					return "~gOpen"
				}
				return "~yPrivate"
			}
			return "~Runknown"
		}())
		count++
	}
	imc_printf(ch, libc.CString("\r\n~W%d ~cchannels found."), count)
	imc_to_char(libc.CString("\r\n~RRed ~clocal name indicates a channel not being listened to.\r\n"), ch)
}
func imclisten(ch *char_data, argument *byte) {
	var c *IMC_CHANNEL
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("~cCurrently tuned into:\r\n"), ch)
		if ch.Player_specials.Imcchardata.Imc_listen != nil && *ch.Player_specials.Imcchardata.Imc_listen != '\x00' {
			imc_printf(ch, libc.CString("~W%s"), ch.Player_specials.Imcchardata.Imc_listen)
		} else {
			imc_to_char(libc.CString("~WNone"), ch)
		}
		imc_to_char(libc.CString("\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("all")) == 0 {
		for c = first_imc_channel; c != nil; c = c.Next {
			if c.Local_name == nil {
				continue
			}
			if ch.Player_specials.Imcchardata.Imcperm >= int(c.Level) && !imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
				imc_addname(&ch.Player_specials.Imcchardata.Imc_listen, c.Local_name)
			}
		}
		imc_to_char(libc.CString("~YYou are now listening to all available IMC2 channels.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("none")) == 0 {
		for c = first_imc_channel; c != nil; c = c.Next {
			if c.Local_name == nil {
				continue
			}
			if imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
				imc_removename(&ch.Player_specials.Imcchardata.Imc_listen, c.Local_name)
			}
		}
		imc_to_char(libc.CString("~YYou no longer listen to any available IMC2 channels.\r\n"), ch)
		return
	}
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(argument)
		return c
	}()) == nil {
		imc_to_char(libc.CString("No such channel configured locally.\r\n"), ch)
		return
	}
	if ch.Player_specials.Imcchardata.Imcperm < int(c.Level) {
		imc_to_char(libc.CString("No such channel configured locally.\r\n"), ch)
		return
	}
	if imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
		imc_removename(&ch.Player_specials.Imcchardata.Imc_listen, c.Local_name)
		imc_to_char(libc.CString("Channel off.\r\n"), ch)
	} else {
		imc_addname(&ch.Player_specials.Imcchardata.Imc_listen, c.Local_name)
		imc_to_char(libc.CString("Channel on.\r\n"), ch)
	}
}
func imctell(ch *char_data, argument *byte) {
	var (
		buf  [4096]byte
		buf1 [4096]byte
	)
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 1)) != 0 {
		imc_to_char(libc.CString("You are not authorized to use imctell.\r\n"), ch)
		return
	}
	argument = imcone_argument(argument, &buf[0])
	if argument == nil || *argument == '\x00' {
		var x int
		imc_to_char(libc.CString("Usage: imctell user@mud <message>\r\n"), ch)
		imc_to_char(libc.CString("Usage: imctell [on]/[off]\r\n\r\n"), ch)
		imc_printf(ch, libc.CString("~cThe last %d things you were told:\r\n"), MAX_IMCTELLHISTORY)
		for x = 0; x < MAX_IMCTELLHISTORY; x++ {
			if (ch.Player_specials.Imcchardata.Imc_tellhistory[x]) == nil {
				break
			}
			imc_to_char(ch.Player_specials.Imcchardata.Imc_tellhistory[x], ch)
		}
		return
	}
	if C.strcasecmp(argument, libc.CString("on")) == 0 {
		ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 0)
		imc_to_char(libc.CString("You now send and receive imctells.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("off")) == 0 {
		ch.Player_specials.Imcchardata.Imcflag |= 1 << 0
		imc_to_char(libc.CString("You no longer send and receive imctells.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 0)) != 0 {
		imc_to_char(libc.CString("You have imctells turned off.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
		imc_to_char(libc.CString("You are invisible.\r\n"), ch)
		return
	}
	if !check_mudof(ch, &buf[0]) {
		return
	}
	if *argument == '@' {
		var (
			p    *byte
			p2   *byte
			buf2 [1024]byte
		)
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		strlcpy(&buf2[0], argument, SMST)
		p = imc_send_social(ch, argument, 1)
		if p == nil || *p == '\x00' {
			return
		}
		imc_send_tell(GET_NAME(ch), &buf[0], p, 2)
		p2 = imc_send_social(ch, &buf2[0], 2)
		if p2 == nil || *p2 == '\x00' {
			return
		}
		stdio.Snprintf(&buf1[0], LGST, "~WImctell ~C%s: ~c%s\r\n", &buf[0], p2)
	} else if *argument == ',' {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		imc_send_tell(GET_NAME(ch), &buf[0], color_mtoi(argument), 1)
		stdio.Snprintf(&buf1[0], LGST, "~WImctell: ~c%s %s\r\n", &buf[0], argument)
	} else {
		imc_send_tell(GET_NAME(ch), &buf[0], color_mtoi(argument), 0)
		stdio.Snprintf(&buf1[0], LGST, "~cYou imctell ~C%s ~c'~W%s~c'\r\n", &buf[0], argument)
	}
	imc_to_char(&buf1[0], ch)
	imc_update_tellhistory(ch, &buf1[0])
}
func imcreply(ch *char_data, argument *byte) {
	var buf1 [4096]byte
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 1)) != 0 {
		imc_to_char(libc.CString("You are not authorized to use imcreply.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 0)) != 0 {
		imc_to_char(libc.CString("You have imctells turned off.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
		imc_to_char(libc.CString("You are invisible.\r\n"), ch)
		return
	}
	if ch.Player_specials.Imcchardata.Rreply == nil {
		imc_to_char(libc.CString("You haven't received an imctell yet.\r\n"), ch)
		return
	}
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("imcreply what?\r\n"), ch)
		return
	}
	if !check_mudof(ch, ch.Player_specials.Imcchardata.Rreply) {
		return
	}
	if *argument == '@' {
		var (
			p    *byte
			p2   *byte
			buf2 [1024]byte
		)
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		strlcpy(&buf2[0], argument, SMST)
		p = imc_send_social(ch, argument, 1)
		if p == nil || *p == '\x00' {
			return
		}
		imc_send_tell(GET_NAME(ch), ch.Player_specials.Imcchardata.Rreply, p, 2)
		p2 = imc_send_social(ch, &buf2[0], 2)
		if p2 == nil || *p2 == '\x00' {
			return
		}
		stdio.Snprintf(&buf1[0], LGST, "~WImctell ~C%s: ~c%s\r\n", ch.Player_specials.Imcchardata.Rreply, p2)
	} else if *argument == ',' {
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		imc_send_tell(GET_NAME(ch), ch.Player_specials.Imcchardata.Rreply, color_mtoi(argument), 1)
		stdio.Snprintf(&buf1[0], LGST, "~WImctell: ~c%s %s\r\n", ch.Player_specials.Imcchardata.Rreply, argument)
	} else {
		imc_send_tell(GET_NAME(ch), ch.Player_specials.Imcchardata.Rreply, color_mtoi(argument), 0)
		stdio.Snprintf(&buf1[0], LGST, "~cYou imctell ~C%s ~c'~W%s~c'\r\n", ch.Player_specials.Imcchardata.Rreply, argument)
	}
	imc_to_char(&buf1[0], ch)
	imc_update_tellhistory(ch, &buf1[0])
}
func imcwho(ch *char_data, argument *byte) {
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("imcwho which mud? See imclist for a list of connected muds.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, this_imcmud.Localname) == 0 && ch.Player_specials.Imcchardata.Imcperm >= IMCPERM_IMM {
		imc_to_char(imc_assemble_who(), ch)
		return
	}
	if !check_mud(ch, argument) {
		return
	}
	imc_send_who(GET_NAME(ch), argument, libc.CString("who"))
}
func imclocate(ch *char_data, argument *byte) {
	var user [1024]byte
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("imclocate who?\r\n"), ch)
		return
	}
	stdio.Snprintf(&user[0], SMST, "%s@*", argument)
	imc_send_whois(GET_NAME(ch), &user[0])
}
func imcfinger(ch *char_data, argument *byte) {
	var (
		name [4096]byte
		arg  [1024]byte
	)
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 6)) != 0 {
		imc_to_char(libc.CString("You are not authorized to use imcfinger.\r\n"), ch)
		return
	}
	argument = imcone_argument(argument, &arg[0])
	if arg[0] == '\x00' {
		imc_to_char(libc.CString("~wUsage: imcfinger person@mud\r\n"), ch)
		imc_to_char(libc.CString("~wUsage: imcfinger <field> <value>\r\n"), ch)
		imc_to_char(libc.CString("~wWhere field is one of:\r\n\r\n"), ch)
		imc_to_char(libc.CString("~wdisplay email homepage icq aim yahoo msn privacy comment\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("display")) == 0 {
		imc_to_char(libc.CString("~GYour current information:\r\n\r\n"), ch)
		imc_printf(ch, libc.CString("~GEmail   : ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Email != nil && *ch.Player_specials.Imcchardata.Email != '\x00' {
				return ch.Player_specials.Imcchardata.Email
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GHomepage: ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Homepage != nil && *ch.Player_specials.Imcchardata.Homepage != '\x00' {
				return ch.Player_specials.Imcchardata.Homepage
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GICQ     : ~g%d\r\n"), ch.Player_specials.Imcchardata.Icq)
		imc_printf(ch, libc.CString("~GAIM     : ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Aim != nil && *ch.Player_specials.Imcchardata.Aim != '\x00' {
				return ch.Player_specials.Imcchardata.Aim
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GYahoo   : ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Yahoo != nil && *ch.Player_specials.Imcchardata.Yahoo != '\x00' {
				return ch.Player_specials.Imcchardata.Yahoo
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GMSN     : ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Msn != nil && *ch.Player_specials.Imcchardata.Msn != '\x00' {
				return ch.Player_specials.Imcchardata.Msn
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GComment : ~g%s\r\n"), func() *byte {
			if ch.Player_specials.Imcchardata.Comment != nil && *ch.Player_specials.Imcchardata.Comment != '\x00' {
				return ch.Player_specials.Imcchardata.Comment
			}
			return libc.CString("None")
		}())
		imc_printf(ch, libc.CString("~GPrivacy : ~g%s\r\n"), func() string {
			if (ch.Player_specials.Imcchardata.Imcflag & (1 << 5)) != 0 {
				return "Enabled"
			}
			return "Disabled"
		}())
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("privacy")) == 0 {
		if (ch.Player_specials.Imcchardata.Imcflag & (1 << 5)) != 0 {
			ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 5)
			imc_to_char(libc.CString("Privacy flag removed. Your information will now be visible on imcfinger.\r\n"), ch)
		} else {
			ch.Player_specials.Imcchardata.Imcflag |= 1 << 5
			imc_to_char(libc.CString("Privacy flag enabled. Your information will no longer be visible on imcfinger.\r\n"), ch)
		}
		return
	}
	if argument == nil || *argument == '\x00' {
		if int(this_imcmud.State) != IMC_ONLINE {
			imc_to_char(libc.CString("The mud is not currently connected to IMC2.\r\n"), ch)
			return
		}
		if !check_mudof(ch, &arg[0]) {
			return
		}
		stdio.Snprintf(&name[0], LGST, "finger %s", imc_nameof(&arg[0]))
		imc_send_who(GET_NAME(ch), imc_mudof(&arg[0]), &name[0])
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("email")) == 0 {
		for {
			if ch.Player_specials.Imcchardata.Email != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Email))
				ch.Player_specials.Imcchardata.Email = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Email = C.strdup(argument)
		imc_printf(ch, libc.CString("Your email address has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Email)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("homepage")) == 0 {
		for {
			if ch.Player_specials.Imcchardata.Homepage != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Homepage))
				ch.Player_specials.Imcchardata.Homepage = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Homepage = C.strdup(argument)
		imc_printf(ch, libc.CString("Your homepage has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Homepage)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("icq")) == 0 {
		ch.Player_specials.Imcchardata.Icq = libc.Atoi(libc.GoString(argument))
		imc_printf(ch, libc.CString("Your ICQ Number has changed to: %d\r\n"), ch.Player_specials.Imcchardata.Icq)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("aim")) == 0 {
		for {
			if ch.Player_specials.Imcchardata.Aim != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Aim))
				ch.Player_specials.Imcchardata.Aim = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Aim = C.strdup(argument)
		imc_printf(ch, libc.CString("Your AIM Screenname has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Aim)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("yahoo")) == 0 {
		for {
			if ch.Player_specials.Imcchardata.Yahoo != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Yahoo))
				ch.Player_specials.Imcchardata.Yahoo = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Yahoo = C.strdup(argument)
		imc_printf(ch, libc.CString("Your Yahoo Screenname has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Yahoo)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("msn")) == 0 {
		for {
			if ch.Player_specials.Imcchardata.Msn != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Msn))
				ch.Player_specials.Imcchardata.Msn = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Msn = C.strdup(argument)
		imc_printf(ch, libc.CString("Your MSN Screenname has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Msn)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("comment")) == 0 {
		if C.strlen(argument) > 78 {
			imc_to_char(libc.CString("You must limit the comment line to 78 characters or less.\r\n"), ch)
			return
		}
		for {
			if ch.Player_specials.Imcchardata.Comment != nil {
				libc.Free(unsafe.Pointer(ch.Player_specials.Imcchardata.Comment))
				ch.Player_specials.Imcchardata.Comment = nil
			}
			if true {
				break
			}
		}
		ch.Player_specials.Imcchardata.Comment = C.strdup(argument)
		imc_printf(ch, libc.CString("Your comment line has changed to: %s\r\n"), ch.Player_specials.Imcchardata.Comment)
		return
	}
	imcfinger(ch, libc.CString(""))
}
func imcinfo(ch *char_data, argument *byte) {
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Syntax: imcinfo <mud>\r\n"), ch)
		return
	}
	if !check_mud(ch, argument) {
		return
	}
	imc_send_who(GET_NAME(ch), argument, libc.CString("info"))
}
func imcbeep(ch *char_data, argument *byte) {
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 3)) != 0 {
		imc_to_char(libc.CString("You are not authorized to use imcbeep.\r\n"), ch)
		return
	}
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Usage: imcbeep user@mud\r\n"), ch)
		imc_to_char(libc.CString("Usage: imcbeep [on]/[off]\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("on")) == 0 {
		ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 2)
		imc_to_char(libc.CString("You now send and receive imcbeeps.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("off")) == 0 {
		ch.Player_specials.Imcchardata.Imcflag |= 1 << 2
		imc_to_char(libc.CString("You no longer send and receive imcbeeps.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 2)) != 0 {
		imc_to_char(libc.CString("You have imcbeep turned off.\r\n"), ch)
		return
	}
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
		imc_to_char(libc.CString("You are invisible.\r\n"), ch)
		return
	}
	if !check_mudof(ch, argument) {
		return
	}
	imc_send_beep(GET_NAME(ch), argument)
	imc_printf(ch, libc.CString("~cYou imcbeep ~Y%s~c.\r\n"), argument)
}
func imclist(ch *char_data, argument *byte) {
	var (
		p          *REMOTEINFO
		serverpath [4096]byte
		netname    [1024]byte
		start      *byte
		onpath     *byte
		count      int = 1
		end        int
	)
	if argument != nil && *argument != '\x00' {
		imcinfo(ch, argument)
		return
	}
	imcpager_printf(ch, libc.CString("~WActive muds on %s:~!\r\n"), this_imcmud.Network)
	imcpager_printf(ch, libc.CString("~c%-15.15s ~B%-40.40s~! ~g%-15.15s ~G%s"), "Name", "IMC2 Version", "Network", "Server")
	imcpager_printf(ch, libc.CString("\r\n\r\n~c%-15.15s ~B%-40.40s ~g%-15.15s ~G%s"), this_imcmud.Localname, this_imcmud.Versionid, this_imcmud.Network, this_imcmud.Servername)
	for p = first_rinfo; p != nil; func() int {
		p = p.Next
		return func() int {
			p := &count
			x := *p
			*p++
			return x
		}()
	}() {
		if C.strcasecmp(p.Network, libc.CString("unknown")) == 0 {
			strlcpy(&netname[0], this_imcmud.Network, SMST)
		} else {
			strlcpy(&netname[0], p.Network, SMST)
		}
		if p.Path != nil && *p.Path != '\x00' {
			if (func() *byte {
				start = C.strchr(p.Path, '!')
				return start
			}()) != nil {
				start = (*byte)(unsafe.Add(unsafe.Pointer(start), 1))
				onpath = start
				end = 0
				for onpath = start; *onpath != '!' && *onpath != '\x00'; onpath = (*byte)(unsafe.Add(unsafe.Pointer(onpath), 1)) {
					serverpath[end] = *onpath
					end++
				}
				serverpath[end] = '\x00'
			} else {
				strlcpy(&serverpath[0], p.Path, LGST)
			}
		}
		imcpager_printf(ch, libc.CString("\r\n~%c%-15.15s ~B%-40.40s ~g%-15.15s ~G%s"), func() int {
			if p.Expired {
				return 'R'
			}
			return 'c'
		}(), p.Name, p.Version, &netname[0], &serverpath[0])
	}
	imcpager_printf(ch, libc.CString("\r\n~WRed mud names indicate connections that are down."))
	imcpager_printf(ch, libc.CString("\r\n~W%d muds on %s found.\r\n"), count, this_imcmud.Network)
}
func imcconnect(ch *char_data, argument *byte) {
	if this_imcmud != nil && int(this_imcmud.State) > IMC_OFFLINE {
		imc_to_char(libc.CString("The IMC2 network connection appears to already be engaged!\r\n"), ch)
		return
	}
	imcconnect_attempts = 0
	imcwait = 0
	imc_startup(TRUE != 0, -1, FALSE != 0)
}
func imcdisconnect(ch *char_data, argument *byte) {
	if this_imcmud != nil && int(this_imcmud.State) == IMC_OFFLINE {
		imc_to_char(libc.CString("The IMC2 network connection does not appear to be engaged!\r\n"), ch)
		return
	}
	imc_shutdown(FALSE != 0)
}
func imcconfig(ch *char_data, argument *byte) {
	var arg1 [1024]byte
	argument = imcone_argument(argument, &arg1[0])
	if arg1[0] == '\x00' {
		imc_to_char(libc.CString("~wSyntax: &Gimc <field> [value]\r\n\r\n"), ch)
		imc_to_char(libc.CString("~wConfiguration info for your mud. Changes save when edited.\r\n"), ch)
		imc_to_char(libc.CString("~wYou may set the following:\r\n\r\n"), ch)
		imc_to_char(libc.CString("~wShow           : ~GDisplays your current configuration.\r\n"), ch)
		imc_to_char(libc.CString("~wLocalname      : ~GThe name IMC2 knows your mud by.\r\n"), ch)
		imc_to_char(libc.CString("~wAutoconnect    : ~GToggles automatic connection on reboots.\r\n"), ch)
		imc_to_char(libc.CString("~wMinPlayerLevel : ~GSets the minimum level IMC2 can see your players at.\r\n"), ch)
		imc_to_char(libc.CString("~wMinImmLevel    : ~GSets the level at which immortal commands become available.\r\n"), ch)
		imc_to_char(libc.CString("~wAdminlevel     : ~GSets the level at which administrative commands become available.\r\n"), ch)
		imc_to_char(libc.CString("~wImplevel       : ~GSets the level at which immplementor commands become available.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoname       : ~GName of your mud, as seen from the imcquery info sheet.\r\n"), ch)
		imc_to_char(libc.CString("~wInfohost       : ~GTelnet address of your mud.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoport       : ~GTelnet port of your mud.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoemail      : ~GEmail address of the mud's IMC administrator.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoWWW        : ~GThe Web address of your mud.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoBase       : ~GThe codebase your mud uses.\r\n"), ch)
		imc_to_char(libc.CString("~wInfoDetails    : ~GSHORT Description of your mud.\r\n"), ch)
		imc_to_char(libc.CString("~wServerAddr     : ~GDNS or IP address of the server you mud connects to.\r\n"), ch)
		imc_to_char(libc.CString("~wServerPort     : ~GPort of the server your mud connects to.\r\n"), ch)
		imc_to_char(libc.CString("~wClientPwd      : ~GClient password for your mud.\r\n"), ch)
		imc_to_char(libc.CString("~wServerPwd      : ~GServer password for your mud.\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("sha256")) == 0 {
		this_imcmud.Sha256 = !this_imcmud.Sha256
		if this_imcmud.Sha256 {
			imc_to_char(libc.CString("SHA-256 support enabled.\r\n"), ch)
		} else {
			imc_to_char(libc.CString("SHA-256 support disabled.\r\n"), ch)
		}
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("sha256pass")) == 0 {
		this_imcmud.Sha256pass = !this_imcmud.Sha256pass
		if this_imcmud.Sha256pass {
			imc_to_char(libc.CString("SHA-256 Authentication enabled.\r\n"), ch)
		} else {
			imc_to_char(libc.CString("SHA-256 Authentication disabled.\r\n"), ch)
		}
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("autoconnect")) == 0 {
		this_imcmud.Autoconnect = !this_imcmud.Autoconnect
		if this_imcmud.Autoconnect {
			imc_to_char(libc.CString("Autoconnect enabled.\r\n"), ch)
		} else {
			imc_to_char(libc.CString("Autoconnect disabled.\r\n"), ch)
		}
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("show")) == 0 {
		imc_printf(ch, libc.CString("~wLocalname      : ~G%s\r\n"), this_imcmud.Localname)
		imc_printf(ch, libc.CString("~wAutoconnect    : ~G%s\r\n"), func() string {
			if this_imcmud.Autoconnect {
				return "Enabled"
			}
			return "Disabled"
		}())
		imc_printf(ch, libc.CString("~wMinPlayerLevel : ~G%d\r\n"), this_imcmud.Minlevel)
		imc_printf(ch, libc.CString("~wMinImmLevel    : ~G%d\r\n"), this_imcmud.Immlevel)
		imc_printf(ch, libc.CString("~wAdminlevel     : ~G%d\r\n"), this_imcmud.Adminlevel)
		imc_printf(ch, libc.CString("~wImplevel       : ~G%d\r\n"), this_imcmud.Implevel)
		imc_printf(ch, libc.CString("~wInfoname       : ~G%s\r\n"), this_imcmud.Fullname)
		imc_printf(ch, libc.CString("~wInfohost       : ~G%s\r\n"), this_imcmud.Ihost)
		imc_printf(ch, libc.CString("~wInfoport       : ~G%d\r\n"), this_imcmud.Iport)
		imc_printf(ch, libc.CString("~wInfoemail      : ~G%s\r\n"), this_imcmud.Email)
		imc_printf(ch, libc.CString("~wInfoWWW        : ~G%s\r\n"), this_imcmud.Www)
		imc_printf(ch, libc.CString("~wInfoBase       : ~G%s\r\n"), this_imcmud.Base)
		imc_printf(ch, libc.CString("~wInfoDetails    : ~G%s\r\n\r\n"), this_imcmud.Details)
		imc_printf(ch, libc.CString("~wServerAddr     : ~G%s\r\n"), this_imcmud.Rhost)
		imc_printf(ch, libc.CString("~wServerPort     : ~G%d\r\n"), this_imcmud.Rport)
		imc_printf(ch, libc.CString("~wClientPwd      : ~G%s\r\n"), this_imcmud.Clientpw)
		imc_printf(ch, libc.CString("~wServerPwd      : ~G%s\r\n"), this_imcmud.Serverpw)
		if this_imcmud.Sha256 {
			imc_to_char(libc.CString("~RThis mud has enabled SHA-256 authentication.\r\n"), ch)
		} else {
			imc_to_char(libc.CString("~RThis mud has disabled SHA-256 authentication.\r\n"), ch)
		}
		if this_imcmud.Sha256 && this_imcmud.Sha256pass {
			imc_to_char(libc.CString("~RThe mud is using SHA-256 encryption to authenticate.\r\n"), ch)
		} else {
			imc_to_char(libc.CString("~RThe mud is using plain text passwords to authenticate.\r\n"), ch)
		}
		return
	}
	if argument == nil || *argument == '\x00' {
		imcconfig(ch, libc.CString(""))
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("minplayerlevel")) == 0 {
		var value int = libc.Atoi(libc.GoString(argument))
		imc_printf(ch, libc.CString("Minimum level set to %d\r\n"), value)
		this_imcmud.Minlevel = value
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("minimmlevel")) == 0 {
		var value int = libc.Atoi(libc.GoString(argument))
		imc_printf(ch, libc.CString("Immortal level set to %d\r\n"), value)
		this_imcmud.Immlevel = value
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("adminlevel")) == 0 {
		var value int = libc.Atoi(libc.GoString(argument))
		imc_printf(ch, libc.CString("Admin level set to %d\r\n"), value)
		this_imcmud.Adminlevel = value
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("implevel")) == 0 && ch.Player_specials.Imcchardata.Imcperm == IMCPERM_IMP {
		var value int = libc.Atoi(libc.GoString(argument))
		imc_printf(ch, libc.CString("Implementor level set to %d\r\n"), value)
		this_imcmud.Implevel = value
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infoname")) == 0 {
		for {
			if this_imcmud.Fullname != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Fullname))
				this_imcmud.Fullname = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Fullname = C.strdup(argument)
		imc_save_config()
		imc_printf(ch, libc.CString("Infoname change to %s\r\n"), argument)
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infohost")) == 0 {
		for {
			if this_imcmud.Ihost != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Ihost))
				this_imcmud.Ihost = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Ihost = C.strdup(argument)
		imc_save_config()
		imc_printf(ch, libc.CString("Infohost changed to %s\r\n"), argument)
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infoport")) == 0 {
		this_imcmud.Iport = libc.Atoi(libc.GoString(argument))
		imc_save_config()
		imc_printf(ch, libc.CString("Infoport changed to %d\r\n"), this_imcmud.Iport)
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infoemail")) == 0 {
		for {
			if this_imcmud.Email != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Email))
				this_imcmud.Email = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Email = C.strdup(argument)
		imc_save_config()
		imc_printf(ch, libc.CString("Infoemail changed to %s\r\n"), argument)
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infowww")) == 0 {
		for {
			if this_imcmud.Www != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Www))
				this_imcmud.Www = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Www = C.strdup(argument)
		imc_save_config()
		imc_printf(ch, libc.CString("InfoWWW changed to %s\r\n"), argument)
		imc_send_keepalive(nil, libc.CString("*@*"))
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infobase")) == 0 {
		var cbase [1024]byte
		for {
			if this_imcmud.Base != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Base))
				this_imcmud.Base = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Base = C.strdup(argument)
		imc_save_config()
		imc_printf(ch, libc.CString("Infobase changed to %s\r\n"), argument)
		for {
			if this_imcmud.Versionid != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Versionid))
				this_imcmud.Versionid = nil
			}
			if true {
				break
			}
		}
		stdio.Snprintf(&cbase[0], SMST, "%s%s", IMC_VERSION_STRING, this_imcmud.Base)
		this_imcmud.Versionid = C.strdup(&cbase[0])
		imc_send_keepalive(nil, libc.CString("*@*"))
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("infodetails")) == 0 {
		for {
			if this_imcmud.Details != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Details))
				this_imcmud.Details = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Details = C.strdup(argument)
		imc_save_config()
		imc_to_char(libc.CString("Infodetails updated.\r\n"), ch)
		return
	}
	if int(this_imcmud.State) != IMC_OFFLINE {
		imc_printf(ch, libc.CString("Cannot alter %s while the mud is connected to IMC.\r\n"), &arg1[0])
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("serveraddr")) == 0 {
		for {
			if this_imcmud.Rhost != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Rhost))
				this_imcmud.Rhost = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Rhost = C.strdup(argument)
		imc_printf(ch, libc.CString("ServerAddr changed to %s\r\n"), argument)
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("serverport")) == 0 {
		this_imcmud.Rport = uint16(int16(libc.Atoi(libc.GoString(argument))))
		imc_printf(ch, libc.CString("ServerPort changed to %d\r\n"), this_imcmud.Rport)
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("clientpwd")) == 0 {
		for {
			if this_imcmud.Clientpw != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Clientpw))
				this_imcmud.Clientpw = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Clientpw = C.strdup(argument)
		imc_printf(ch, libc.CString("Clientpwd changed to %s\r\n"), argument)
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("serverpwd")) == 0 {
		for {
			if this_imcmud.Serverpw != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Serverpw))
				this_imcmud.Serverpw = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Serverpw = C.strdup(argument)
		imc_printf(ch, libc.CString("Serverpwd changed to %s\r\n"), argument)
		imc_save_config()
		return
	}
	if C.strcasecmp(&arg1[0], libc.CString("localname")) == 0 {
		for {
			if this_imcmud.Localname != nil {
				libc.Free(unsafe.Pointer(this_imcmud.Localname))
				this_imcmud.Localname = nil
			}
			if true {
				break
			}
		}
		this_imcmud.Localname = C.strdup(argument)
		this_imcmud.Sha256pass = FALSE != 0
		imc_save_config()
		imc_printf(ch, libc.CString("Localname changed to %s\r\n"), argument)
		return
	}
	imcconfig(ch, libc.CString(""))
}
func imcignore(ch *char_data, argument *byte) {
	var (
		count int
		ign   *IMC_IGNORE
		arg   [1024]byte
	)
	argument = imcone_argument(argument, &arg[0])
	imc_to_char(libc.CString("IMCignore is disabled at this time.\r\n"), ch)
	return
	if arg[0] == '\x00' {
		imc_to_char(libc.CString("You currently ignore the following:\r\n"), ch)
		for func() *IMC_IGNORE {
			count = 0
			return func() *IMC_IGNORE {
				ign = ch.Player_specials.Imcchardata.Imcfirst_ignore
				return ign
			}()
		}(); ign != nil; func() int {
			ign = ign.Next
			return func() int {
				p := &count
				x := *p
				*p++
				return x
			}()
		}() {
			imc_printf(ch, libc.CString("%s\r\n"), ign.Name)
		}
		if count == 0 {
			imc_to_char(libc.CString(" none\r\n"), ch)
		} else {
			imc_printf(ch, libc.CString("\r\n[total %d]\r\n"), count)
		}
		imc_to_char(libc.CString("For help on imcignore, type: IMCIGNORE HELP\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("help")) == 0 {
		imc_to_char(libc.CString("~wTo see your current ignores  : ~GIMCIGNORE\r\n"), ch)
		imc_to_char(libc.CString("~wTo add an ignore             : ~GIMCIGNORE ADD <argument>\r\n"), ch)
		imc_to_char(libc.CString("~wTo delete an ignore          : ~GIMCIGNORE DELETE <argument>\r\n"), ch)
		imc_to_char(libc.CString("~WSee your MUD's help for more information.\r\n"), ch)
		return
	}
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Must specify both action and name.\r\n"), ch)
		imc_to_char(libc.CString("Please see IMCIGNORE HELP for details.\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("delete")) == 0 {
		for ign = ch.Player_specials.Imcchardata.Imcfirst_ignore; ign != nil; ign = ign.Next {
			if C.strcasecmp(ign.Name, argument) == 0 {
				for {
					if ign.Prev == nil {
						ch.Player_specials.Imcchardata.Imcfirst_ignore = ign.Next
						if ch.Player_specials.Imcchardata.Imcfirst_ignore != nil {
							ch.Player_specials.Imcchardata.Imcfirst_ignore.Prev = nil
						}
					} else {
						ign.Prev.Next = ign.Next
					}
					if ign.Next == nil {
						ch.Player_specials.Imcchardata.Imclast_ignore = ign.Prev
						if ch.Player_specials.Imcchardata.Imclast_ignore != nil {
							ch.Player_specials.Imcchardata.Imclast_ignore.Next = nil
						}
					} else {
						ign.Next.Prev = ign.Prev
					}
					if true {
						break
					}
				}
				for {
					if ign.Name != nil {
						libc.Free(unsafe.Pointer(ign.Name))
						ign.Name = nil
					}
					if true {
						break
					}
				}
				for {
					if ign != nil {
						libc.Free(unsafe.Pointer(ign))
						ign = nil
					}
					if true {
						break
					}
				}
				imc_to_char(libc.CString("Entry deleted.\r\n"), ch)
				return
			}
		}
		imc_to_char(libc.CString("Entry not found.\r\nPlease check your ignores by typing IMCIGNORE with no arguments.\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("add")) == 0 {
		for {
			if (func() *IMC_IGNORE {
				ign = new(IMC_IGNORE)
				return ign
			}()) == nil {
				imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
				abort()
			}
			if true {
				break
			}
		}
		ign.Name = C.strdup(argument)
		for {
			if ch.Player_specials.Imcchardata.Imcfirst_ignore == nil {
				ch.Player_specials.Imcchardata.Imcfirst_ignore = ign
				ch.Player_specials.Imcchardata.Imclast_ignore = ign
			} else {
				ch.Player_specials.Imcchardata.Imclast_ignore.Next = ign
			}
			ign.Next = nil
			if ch.Player_specials.Imcchardata.Imcfirst_ignore == ign {
				ign.Prev = nil
			} else {
				ign.Prev = ch.Player_specials.Imcchardata.Imclast_ignore
			}
			ch.Player_specials.Imcchardata.Imclast_ignore = ign
			if true {
				break
			}
		}
		imc_printf(ch, libc.CString("%s will now be ignored.\r\n"), argument)
		return
	}
	imcignore(ch, libc.CString("help"))
}
func imcban(ch *char_data, argument *byte) {
	var (
		count int
		ban   *IMC_BAN
		arg   [1024]byte
	)
	argument = imcone_argument(argument, &arg[0])
	if arg[0] == '\x00' {
		imc_to_char(libc.CString("The mud currently bans the following:\r\n"), ch)
		for func() *IMC_BAN {
			count = 0
			return func() *IMC_BAN {
				ban = first_imc_ban
				return ban
			}()
		}(); ban != nil; func() int {
			ban = ban.Next
			return func() int {
				p := &count
				x := *p
				*p++
				return x
			}()
		}() {
			imc_printf(ch, libc.CString("%s\r\n"), ban.Name)
		}
		if count == 0 {
			imc_to_char(libc.CString(" none\r\n"), ch)
		} else {
			imc_printf(ch, libc.CString("\r\n[total %d]\r\n"), count)
		}
		imc_to_char(libc.CString("Type: IMCBAN HELP for more information.\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("help")) == 0 {
		imc_to_char(libc.CString("~wTo see the current bans             : ~GIMCBAN\r\n"), ch)
		imc_to_char(libc.CString("~wTo add a MUD to the ban list        : ~GIMCBAN ADD <argument>\r\n"), ch)
		imc_to_char(libc.CString("~wTo delete a MUD from the ban list   : ~GIMCBAN DELETE <argument>\r\n"), ch)
		imc_to_char(libc.CString("~WSee your MUD's help for more information.\r\n"), ch)
		return
	}
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Must specify both action and name.\r\nPlease type IMCBAN HELP for more information\r\n"), ch)
		return
	}
	if C.strcasecmp(&arg[0], libc.CString("delete")) == 0 {
		if imc_delban(argument) {
			imc_savebans()
			imc_to_char(libc.CString("Entry deleted.\r\n"), ch)
			return
		}
		imc_to_char(libc.CString("Entry not found.\r\nPlease type IMCBAN without arguments to see the current ban list.\r\n"), ch)
	}
	if C.strcasecmp(&arg[0], libc.CString("add")) == 0 {
		imc_addban(argument)
		imc_savebans()
		imc_printf(ch, libc.CString("Mud %s will now be banned.\r\n"), argument)
		return
	}
	imcban(ch, libc.CString(""))
}
func imc_deny_channel(ch *char_data, argument *byte) {
	var (
		vic_name [1024]byte
		victim   *char_data
		channel  *IMC_CHANNEL
	)
	argument = imcone_argument(argument, &vic_name[0])
	if vic_name[0] == '\x00' || argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Usage: imcdeny <person> <local channel name>\r\n"), ch)
		imc_to_char(libc.CString("Usage: imcdeny <person> [tell/beep/finger]\r\n"), ch)
		return
	}
	if (func() *char_data {
		victim = imc_find_user(&vic_name[0])
		return victim
	}()) == nil {
		imc_to_char(libc.CString("No such person is currently online.\r\n"), ch)
		return
	}
	if ch.Player_specials.Imcchardata.Imcperm <= victim.Player_specials.Imcchardata.Imcperm {
		imc_to_char(libc.CString("You cannot alter their settings.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("tell")) == 0 {
		if (victim.Player_specials.Imcchardata.Imcflag & (1 << 1)) == 0 {
			victim.Player_specials.Imcchardata.Imcflag |= 1 << 1
			imc_printf(ch, libc.CString("%s can no longer use imctells.\r\n"), GET_NAME(victim))
			return
		}
		victim.Player_specials.Imcchardata.Imcflag &= ^(1 << 1)
		imc_printf(ch, libc.CString("%s can use imctells again.\r\n"), GET_NAME(victim))
		return
	}
	if C.strcasecmp(argument, libc.CString("beep")) == 0 {
		if (victim.Player_specials.Imcchardata.Imcflag & (1 << 3)) == 0 {
			victim.Player_specials.Imcchardata.Imcflag |= 1 << 3
			imc_printf(ch, libc.CString("%s can no longer use imcbeeps.\r\n"), GET_NAME(victim))
			return
		}
		victim.Player_specials.Imcchardata.Imcflag &= ^(1 << 3)
		imc_printf(ch, libc.CString("%s can use imcbeeps again.\r\n"), GET_NAME(victim))
		return
	}
	if C.strcasecmp(argument, libc.CString("finger")) == 0 {
		if (victim.Player_specials.Imcchardata.Imcflag & (1 << 6)) == 0 {
			victim.Player_specials.Imcchardata.Imcflag |= 1 << 6
			imc_printf(ch, libc.CString("%s can no longer use imcfingers.\r\n"), GET_NAME(victim))
			return
		}
		victim.Player_specials.Imcchardata.Imcflag &= ^(1 << 6)
		imc_printf(ch, libc.CString("%s can use imcfingers again.\r\n"), GET_NAME(victim))
		return
	}
	if (func() *IMC_CHANNEL {
		channel = imc_findchannel(argument)
		return channel
	}()) == nil {
		imc_to_char(libc.CString("Unknown or unconfigured local channel. Check your channel name.\r\n"), ch)
		return
	}
	if imc_hasname(victim.Player_specials.Imcchardata.Imc_denied, channel.Local_name) {
		imc_printf(ch, libc.CString("%s can now listen to %s\r\n"), GET_NAME(victim), channel.Local_name)
		imc_removename(&victim.Player_specials.Imcchardata.Imc_denied, channel.Local_name)
	} else {
		imc_printf(ch, libc.CString("%s can no longer listen to %s\r\n"), GET_NAME(victim), channel.Local_name)
		imc_addname(&victim.Player_specials.Imcchardata.Imc_denied, channel.Local_name)
	}
}
func imcpermstats(ch *char_data, argument *byte) {
	var victim *char_data
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Usage: imcperms <user>\r\n"), ch)
		return
	}
	if (func() *char_data {
		victim = imc_find_user(argument)
		return victim
	}()) == nil {
		imc_to_char(libc.CString("No such person is currently online.\r\n"), ch)
		return
	}
	if victim.Player_specials.Imcchardata.Imcperm < 0 || victim.Player_specials.Imcchardata.Imcperm > IMCPERM_IMP {
		imc_printf(ch, libc.CString("%s has an invalid permission setting!\r\n"), GET_NAME(victim))
		return
	}
	imc_printf(ch, libc.CString("~GPermissions for %s: %s\r\n"), GET_NAME(victim), imcperm_names[victim.Player_specials.Imcchardata.Imcperm])
	imc_printf(ch, libc.CString("~gThese permissions were obtained %s.\r\n"), func() string {
		if (victim.Player_specials.Imcchardata.Imcflag & (1 << 9)) != 0 {
			return "manually via imcpermset"
		}
		return "automatically by level"
	}())
}
func imcpermset(ch *char_data, argument *byte) {
	var (
		victim    *char_data
		arg       [1024]byte
		permvalue int
	)
	argument = imcone_argument(argument, &arg[0])
	if arg[0] == '\x00' {
		imc_to_char(libc.CString("Usage: imcpermset <user> <permission>\r\n"), ch)
		imc_to_char(libc.CString("Permission can be one of: None, Mort, Imm, Admin, Imp\r\n"), ch)
		return
	}
	if (func() *char_data {
		victim = imc_find_user(&arg[0])
		return victim
	}()) == nil {
		imc_to_char(libc.CString("No such person is currently online.\r\n"), ch)
		return
	}
	if C.strcasecmp(argument, libc.CString("override")) == 0 {
		permvalue = -1
	} else {
		permvalue = get_imcpermvalue(argument)
		if !imccheck_permissions(ch, permvalue, victim.Player_specials.Imcchardata.Imcperm, TRUE != 0) {
			return
		}
	}
	if victim.Player_specials.Imcchardata.Imcperm == permvalue {
		imc_printf(ch, libc.CString("%s already has a permission level of %s.\r\n"), GET_NAME(victim), imcperm_names[permvalue])
		return
	}
	if permvalue == -1 {
		victim.Player_specials.Imcchardata.Imcflag &= ^(1 << 9)
		imc_printf(ch, libc.CString("~YPermission flag override has been removed from %s\r\n"), GET_NAME(victim))
		return
	}
	victim.Player_specials.Imcchardata.Imcperm = permvalue
	victim.Player_specials.Imcchardata.Imcflag |= 1 << 9
	imc_printf(ch, libc.CString("~YPermission level for %s has been changed to %s\r\n"), GET_NAME(victim), imcperm_names[permvalue])
	if victim.Player_specials.Imcchardata.Imc_listen != nil && int(this_imcmud.State) == IMC_ONLINE {
		var (
			channel  *IMC_CHANNEL = nil
			channels *byte        = victim.Player_specials.Imcchardata.Imc_listen
		)
		for {
			if *channels == '\x00' {
				break
			}
			channels = imcone_argument(channels, &arg[0])
			if (func() *IMC_CHANNEL {
				channel = imc_findchannel(&arg[0])
				return channel
			}()) == nil {
				imc_removename(&victim.Player_specials.Imcchardata.Imc_listen, &arg[0])
			}
			if channel != nil && victim.Player_specials.Imcchardata.Imcperm < int(channel.Level) {
				imc_removename(&victim.Player_specials.Imcchardata.Imc_listen, &arg[0])
				imc_printf(ch, libc.CString("~WRemoving '%s' level channel: '%s', exceeding new permission of '%s'\r\n"), imcperm_names[channel.Level], channel.Local_name, imcperm_names[victim.Player_specials.Imcchardata.Imcperm])
			}
		}
	}
}
func imcinvis(ch *char_data, argument *byte) {
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 4)) != 0 {
		ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 4)
		imc_to_char(libc.CString("You are now imcvisible.\r\n"), ch)
	} else {
		ch.Player_specials.Imcchardata.Imcflag |= 1 << 4
		imc_to_char(libc.CString("You are now imcinvisible.\r\n"), ch)
	}
}
func imcchanwho(ch *char_data, argument *byte) {
	var (
		c     *IMC_CHANNEL
		p     *IMC_PACKET
		chan_ [1024]byte
		mud   [1024]byte
	)
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Usage: imcchanwho <channel> [<mud> <mud> <mud> <...>|<all>]\r\n"), ch)
		return
	}
	argument = imcone_argument(argument, &chan_[0])
	if (func() *IMC_CHANNEL {
		c = imc_findchannel(&chan_[0])
		return c
	}()) == nil {
		imc_to_char(libc.CString("No such channel.\r\n"), ch)
		return
	}
	if ch.Player_specials.Imcchardata.Imcperm < int(c.Level) {
		imc_to_char(libc.CString("No such channel.\r\n"), ch)
		return
	}
	if !c.Refreshed {
		imc_printf(ch, libc.CString("%s has not been refreshed yet.\r\n"), c.Name)
		return
	}
	if C.strcasecmp(argument, libc.CString("all")) != 0 {
		for *argument != '\x00' {
			argument = imcone_argument(argument, &mud[0])
			if !check_mud(ch, &mud[0]) {
				continue
			}
			p = imc_newpacket(GET_NAME(ch), libc.CString("ice-chan-who"), &mud[0])
			imc_addtopacket(p, libc.CString("level=%d"), ch.Player_specials.Imcchardata.Imcperm)
			imc_addtopacket(p, libc.CString("channel=%s"), c.Name)
			imc_addtopacket(p, libc.CString("lname=%s"), func() *byte {
				if c.Local_name != nil {
					return c.Local_name
				}
				return c.Name
			}())
			imc_write_packet(p)
		}
		return
	}
	p = imc_newpacket(GET_NAME(ch), libc.CString("ice-chan-who"), libc.CString("*"))
	imc_addtopacket(p, libc.CString("level=%d"), ch.Player_specials.Imcchardata.Imcperm)
	imc_addtopacket(p, libc.CString("channel=%s"), c.Name)
	imc_addtopacket(p, libc.CString("lname=%s"), func() *byte {
		if c.Local_name != nil {
			return c.Local_name
		}
		return c.Name
	}())
	imc_write_packet(p)
	imc_printf(ch, libc.CString("~G%s"), get_local_chanwho(c))
}
func imcremoteadmin(ch *char_data, argument *byte) {
	var (
		r      *REMOTEINFO
		server [1024]byte
		cmd    [1024]byte
		to     [1024]byte
		pwd    [4096]byte
		p      *IMC_PACKET
	)
	argument = imcone_argument(argument, &server[0])
	argument = imcone_argument(argument, &pwd[0])
	argument = imcone_argument(argument, &cmd[0])
	if server[0] == '\x00' || cmd[0] == '\x00' {
		imc_to_char(libc.CString("Syntax: imcadmin <server> <password> <command> [<data..>]\r\n"), ch)
		imc_to_char(libc.CString("You must be an approved server administrator to use remote commands.\r\n"), ch)
		return
	}
	if (func() *REMOTEINFO {
		r = imc_find_reminfo(&server[0])
		return r
	}()) == nil {
		imc_printf(ch, libc.CString("~W%s ~cis not a valid mud name.\r\n"), &server[0])
		return
	}
	if r.Expired {
		imc_printf(ch, libc.CString("~W%s ~cis not connected right now.\r\n"), r.Name)
		return
	}
	stdio.Snprintf(&to[0], SMST, "IMC@%s", r.Name)
	p = imc_newpacket(GET_NAME(ch), libc.CString("remote-admin"), &to[0])
	imc_addtopacket(p, libc.CString("command=%s"), &cmd[0])
	if argument != nil && *argument != '\x00' {
		imc_addtopacket(p, libc.CString("data=%s"), argument)
	}
	if this_imcmud.Sha256pass {
		var (
			cryptpw [4096]byte
			hash    *byte
		)
		stdio.Snprintf(&cryptpw[0], LGST, "%ld%s", imc_sequencenumber+1, &pwd[0])
		// todo: replace with hash
		hash = &cryptpw[0]
		imc_addtopacket(p, libc.CString("hash=%s"), hash)
	}
	imc_write_packet(p)
	imc_to_char(libc.CString("Remote command sent.\r\n"), ch)
}
func imchelp(ch *char_data, argument *byte) {
	var (
		buf  [4096]byte
		help *IMC_HELP_DATA
		col  int
		perm int
	)
	if argument == nil || *argument == '\x00' {
		strlcpy(&buf[0], libc.CString("~gHelp is available for the following commands:\r\n"), LGST)
		imcstrlcat(&buf[0], libc.CString("~G---------------------------------------------\r\n"), LGST)
		for perm = IMCPERM_MORT; perm <= ch.Player_specials.Imcchardata.Imcperm; perm++ {
			col = 0
			stdio.Snprintf(&buf[C.strlen(&buf[0])], int(LGST-C.strlen(&buf[0])), "\r\n~g%s helps:~G\r\n", imcperm_names[perm])
			for help = first_imc_help; help != nil; help = help.Next {
				if help.Level != perm {
					continue
				}
				stdio.Snprintf(&buf[C.strlen(&buf[0])], int(LGST-C.strlen(&buf[0])), "%-15s", help.Name)
				if func() int {
					p := &col
					*p++
					return *p
				}()%6 == 0 {
					imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
				}
			}
			if col%6 != 0 {
				imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
			}
		}
		imc_to_pager(&buf[0], ch)
		return
	}
	for help = first_imc_help; help != nil; help = help.Next {
		if C.strcasecmp(help.Name, argument) == 0 {
			if help.Text == nil || *help.Text == '\x00' {
				imc_printf(ch, libc.CString("~gNo inforation available for topic ~W%s~g.\r\n"), help.Name)
			} else {
				imc_printf(ch, libc.CString("~g%s\r\n"), help.Text)
			}
			return
		}
	}
	imc_printf(ch, libc.CString("~gNo help exists for topic ~W%s~g.\r\n"), argument)
}
func imccolor(ch *char_data, argument *byte) {
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 8)) != 0 {
		ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 8)
		imc_to_char(libc.CString("IMC2 color is now off.\r\n"), ch)
	} else {
		ch.Player_specials.Imcchardata.Imcflag |= 1 << 8
		imc_to_char(libc.CString("~RIMC2 c~Yo~Gl~Bo~Pr ~Ris now on. Enjoy :)\r\n"), ch)
	}
}
func imcafk(ch *char_data, argument *byte) {
	if (ch.Player_specials.Imcchardata.Imcflag & (1 << 7)) != 0 {
		ch.Player_specials.Imcchardata.Imcflag &= ^(1 << 7)
		imc_to_char(libc.CString("You are no longer AFK to IMC2.\r\n"), ch)
	} else {
		ch.Player_specials.Imcchardata.Imcflag |= 1 << 7
		imc_to_char(libc.CString("You are now AFK to IMC2.\r\n"), ch)
	}
}
func imcdebug(ch *char_data, argument *byte) {
	imcpacketdebug = !imcpacketdebug
	if imcpacketdebug {
		imc_to_char(libc.CString("Packet debug enabled.\r\n"), ch)
	} else {
		imc_to_char(libc.CString("Packet debug disabled.\r\n"), ch)
	}
}
func imc_show_ucache_contents(ch *char_data, argument *byte) {
	var (
		user  *IMCUCACHE_DATA
		users int = 0
	)
	imc_to_pager(libc.CString("Cached user information\r\n"), ch)
	imc_to_pager(libc.CString("User                          | Gender ( 0 = Male, 1 = Female, 2 = Other )\r\n"), ch)
	imc_to_pager(libc.CString("--------------------------------------------------------------------------\r\n"), ch)
	for user = first_imcucache; user != nil; user = user.Next {
		imcpager_printf(ch, libc.CString("%-30s %d\r\n"), user.Name, user.Gender)
		users++
	}
	imcpager_printf(ch, libc.CString("%d users being cached.\r\n"), users)
}
func imccedit(ch *char_data, argument *byte) {
	var (
		cmd        *IMC_CMD_DATA
		tmp        *IMC_CMD_DATA
		alias      *IMC_ALIAS
		alias_next *IMC_ALIAS
		name       [1024]byte
		option     [1024]byte
		found      bool = FALSE != 0
		aliasfound bool = FALSE != 0
	)
	argument = imcone_argument(argument, &name[0])
	argument = imcone_argument(argument, &option[0])
	if name[0] == '\x00' || option[0] == '\x00' {
		imc_to_char(libc.CString("Usage: imccedit <command> <create|delete|alias|rename|code|permission|connected> <field>.\r\n"), ch)
		return
	}
	for cmd = first_imc_command; cmd != nil; cmd = cmd.Next {
		if C.strcasecmp(cmd.Name, &name[0]) == 0 {
			found = TRUE != 0
			break
		}
		for alias = cmd.First_alias; alias != nil; alias = alias.Next {
			if C.strcasecmp(alias.Name, &name[0]) == 0 {
				aliasfound = TRUE != 0
			}
		}
	}
	if C.strcasecmp(&option[0], libc.CString("create")) == 0 {
		if found {
			imc_printf(ch, libc.CString("~gA command named ~W%s ~galready exists.\r\n"), &name[0])
			return
		}
		if aliasfound {
			imc_printf(ch, libc.CString("~g%s already exists as an alias for another command.\r\n"), &name[0])
			return
		}
		for {
			if (func() *IMC_CMD_DATA {
				cmd = new(IMC_CMD_DATA)
				return cmd
			}()) == nil {
				imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
				abort()
			}
			if true {
				break
			}
		}
		cmd.Name = C.strdup(&name[0])
		cmd.Level = ch.Player_specials.Imcchardata.Imcperm
		cmd.Connected = FALSE != 0
		imc_printf(ch, libc.CString("~gCommand ~W%s ~gcreated.\r\n"), cmd.Name)
		if argument != nil && *argument != '\x00' {
			cmd.Function = imc_function(argument)
			if cmd.Function == nil {
				imc_printf(ch, libc.CString("~gFunction ~W%s ~gdoes not exist - set to NULL.\r\n"), argument)
			}
		} else {
			imc_to_char(libc.CString("~gFunction set to NULL.\r\n"), ch)
			cmd.Function = nil
		}
		for {
			if first_imc_command == nil {
				first_imc_command = cmd
				last_imc_command = cmd
			} else {
				last_imc_command.Next = cmd
			}
			cmd.Next = nil
			if first_imc_command == cmd {
				cmd.Prev = nil
			} else {
				cmd.Prev = last_imc_command
			}
			last_imc_command = cmd
			if true {
				break
			}
		}
		imc_savecommands()
		return
	}
	if !found {
		imc_printf(ch, libc.CString("~gNo command named ~W%s ~gexists.\r\n"), &name[0])
		return
	}
	if !imccheck_permissions(ch, cmd.Level, cmd.Level, FALSE != 0) {
		return
	}
	if C.strcasecmp(&option[0], libc.CString("delete")) == 0 {
		imc_printf(ch, libc.CString("~gCommand ~W%s ~ghas been deleted.\r\n"), cmd.Name)
		for alias = cmd.First_alias; alias != nil; alias = alias_next {
			alias_next = alias.Next
			for {
				if alias.Prev == nil {
					cmd.First_alias = alias.Next
					if cmd.First_alias != nil {
						cmd.First_alias.Prev = nil
					}
				} else {
					alias.Prev.Next = alias.Next
				}
				if alias.Next == nil {
					cmd.Last_alias = alias.Prev
					if cmd.Last_alias != nil {
						cmd.Last_alias.Next = nil
					}
				} else {
					alias.Next.Prev = alias.Prev
				}
				if true {
					break
				}
			}
			for {
				if alias.Name != nil {
					libc.Free(unsafe.Pointer(alias.Name))
					alias.Name = nil
				}
				if true {
					break
				}
			}
			for {
				if alias != nil {
					libc.Free(unsafe.Pointer(alias))
					alias = nil
				}
				if true {
					break
				}
			}
		}
		for {
			if cmd.Prev == nil {
				first_imc_command = cmd.Next
				if first_imc_command != nil {
					first_imc_command.Prev = nil
				}
			} else {
				cmd.Prev.Next = cmd.Next
			}
			if cmd.Next == nil {
				last_imc_command = cmd.Prev
				if last_imc_command != nil {
					last_imc_command.Next = nil
				}
			} else {
				cmd.Next.Prev = cmd.Prev
			}
			if true {
				break
			}
		}
		for {
			if cmd.Name != nil {
				libc.Free(unsafe.Pointer(cmd.Name))
				cmd.Name = nil
			}
			if true {
				break
			}
		}
		for {
			if cmd != nil {
				libc.Free(unsafe.Pointer(cmd))
				cmd = nil
			}
			if true {
				break
			}
		}
		imc_savecommands()
		return
	}
	if C.strcasecmp(&option[0], libc.CString("alias")) == 0 {
		for alias = cmd.First_alias; alias != nil; alias = alias_next {
			alias_next = alias.Next
			if C.strcasecmp(alias.Name, argument) == 0 {
				imc_printf(ch, libc.CString("~W%s ~ghas been removed as an alias for ~W%s\r\n"), argument, cmd.Name)
				for {
					if alias.Prev == nil {
						cmd.First_alias = alias.Next
						if cmd.First_alias != nil {
							cmd.First_alias.Prev = nil
						}
					} else {
						alias.Prev.Next = alias.Next
					}
					if alias.Next == nil {
						cmd.Last_alias = alias.Prev
						if cmd.Last_alias != nil {
							cmd.Last_alias.Next = nil
						}
					} else {
						alias.Next.Prev = alias.Prev
					}
					if true {
						break
					}
				}
				for {
					if alias.Name != nil {
						libc.Free(unsafe.Pointer(alias.Name))
						alias.Name = nil
					}
					if true {
						break
					}
				}
				for {
					if alias != nil {
						libc.Free(unsafe.Pointer(alias))
						alias = nil
					}
					if true {
						break
					}
				}
				imc_savecommands()
				return
			}
		}
		for tmp = first_imc_command; tmp != nil; tmp = tmp.Next {
			if C.strcasecmp(tmp.Name, argument) == 0 {
				imc_printf(ch, libc.CString("~W%s &gis already a command name.\r\n"), argument)
				return
			}
			for alias = tmp.First_alias; alias != nil; alias = alias.Next {
				if C.strcasecmp(argument, alias.Name) == 0 {
					imc_printf(ch, libc.CString("~W%s ~gis already an alias for ~W%s\r\n"), argument, tmp.Name)
					return
				}
			}
		}
		for {
			if (func() *IMC_ALIAS {
				alias = new(IMC_ALIAS)
				return alias
			}()) == nil {
				imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
				abort()
			}
			if true {
				break
			}
		}
		alias.Name = C.strdup(argument)
		for {
			if cmd.First_alias == nil {
				cmd.First_alias = alias
				cmd.Last_alias = alias
			} else {
				cmd.Last_alias.Next = alias
			}
			alias.Next = nil
			if cmd.First_alias == alias {
				alias.Prev = nil
			} else {
				alias.Prev = cmd.Last_alias
			}
			cmd.Last_alias = alias
			if true {
				break
			}
		}
		imc_printf(ch, libc.CString("~W%s ~ghas been added as an alias for ~W%s\r\n"), alias.Name, cmd.Name)
		imc_savecommands()
		return
	}
	if C.strcasecmp(&option[0], libc.CString("connected")) == 0 {
		cmd.Connected = !cmd.Connected
		if cmd.Connected {
			imc_printf(ch, libc.CString("~gCommand ~W%s ~gwill now require a connection to IMC2 to use.\r\n"), cmd.Name)
		} else {
			imc_printf(ch, libc.CString("~gCommand ~W%s ~gwill no longer require a connection to IMC2 to use.\r\n"), cmd.Name)
		}
		imc_savecommands()
		return
	}
	if C.strcasecmp(&option[0], libc.CString("show")) == 0 {
		var buf [4096]byte
		imc_printf(ch, libc.CString("~gCommand       : ~W%s\r\n"), cmd.Name)
		imc_printf(ch, libc.CString("~gPermission    : ~W%s\r\n"), imcperm_names[cmd.Level])
		imc_printf(ch, libc.CString("~gFunction      : ~W%s\r\n"), imc_funcname(cmd.Function))
		imc_printf(ch, libc.CString("~gConnection Req: ~W%s\r\n"), func() string {
			if cmd.Connected {
				return "Yes"
			}
			return "No"
		}())
		if cmd.First_alias != nil {
			var col int = 0
			strlcpy(&buf[0], libc.CString("~gAliases       : ~W"), LGST)
			for alias = cmd.First_alias; alias != nil; alias = alias.Next {
				stdio.Snprintf(&buf[C.strlen(&buf[0])], int(LGST-C.strlen(&buf[0])), "%s ", alias.Name)
				if func() int {
					p := &col
					*p++
					return *p
				}()%10 == 0 {
					imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
				}
			}
			if col%10 != 0 {
				imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
			}
			imc_to_char(&buf[0], ch)
		}
		return
	}
	if argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Required argument missing.\r\n"), ch)
		imccedit(ch, libc.CString(""))
		return
	}
	if C.strcasecmp(&option[0], libc.CString("rename")) == 0 {
		imc_printf(ch, libc.CString("~gCommand ~W%s ~ghas been renamed to ~W%s.\r\n"), cmd.Name, argument)
		for {
			if cmd.Name != nil {
				libc.Free(unsafe.Pointer(cmd.Name))
				cmd.Name = nil
			}
			if true {
				break
			}
		}
		cmd.Name = C.strdup(argument)
		imc_savecommands()
		return
	}
	if C.strcasecmp(&option[0], libc.CString("code")) == 0 {
		cmd.Function = imc_function(argument)
		if cmd.Function == nil {
			imc_printf(ch, libc.CString("~gFunction ~W%s ~gdoes not exist - set to NULL.\r\n"), argument)
		} else {
			imc_printf(ch, libc.CString("~gFunction set to ~W%s.\r\n"), argument)
		}
		imc_savecommands()
		return
	}
	if C.strcasecmp(&option[0], libc.CString("perm")) == 0 || C.strcasecmp(&option[0], libc.CString("permission")) == 0 {
		var permvalue int = get_imcpermvalue(argument)
		if !imccheck_permissions(ch, permvalue, cmd.Level, FALSE != 0) {
			return
		}
		cmd.Level = permvalue
		imc_printf(ch, libc.CString("~gCommand ~W%s ~gpermission level has been changed to ~W%s.\r\n"), cmd.Name, imcperm_names[permvalue])
		imc_savecommands()
		return
	}
	imccedit(ch, libc.CString(""))
}
func imchedit(ch *char_data, argument *byte) {
	var (
		help  *IMC_HELP_DATA
		name  [1024]byte
		cmd   [1024]byte
		found bool = FALSE != 0
	)
	argument = imcone_argument(argument, &name[0])
	argument = imcone_argument(argument, &cmd[0])
	if name[0] == '\x00' || cmd[0] == '\x00' || argument == nil || *argument == '\x00' {
		imc_to_char(libc.CString("Usage: imchedit <topic> [name|perm] <field>\r\n"), ch)
		imc_to_char(libc.CString("Where <field> can be either name, or permission level.\r\n"), ch)
		return
	}
	for help = first_imc_help; help != nil; help = help.Next {
		if C.strcasecmp(help.Name, &name[0]) == 0 {
			found = TRUE != 0
			break
		}
	}
	if !found {
		imc_printf(ch, libc.CString("~gNo help exists for topic ~W%s~g. You will need to add it to the helpfile manually.\r\n"), &name[0])
		return
	}
	if C.strcasecmp(&cmd[0], libc.CString("name")) == 0 {
		imc_printf(ch, libc.CString("~W%s ~ghas been renamed to ~W%s.\r\n"), help.Name, argument)
		for {
			if help.Name != nil {
				libc.Free(unsafe.Pointer(help.Name))
				help.Name = nil
			}
			if true {
				break
			}
		}
		help.Name = C.strdup(argument)
		imc_savehelps()
		return
	}
	if C.strcasecmp(&cmd[0], libc.CString("perm")) == 0 {
		var permvalue int = get_imcpermvalue(argument)
		if !imccheck_permissions(ch, permvalue, help.Level, FALSE != 0) {
			return
		}
		imc_printf(ch, libc.CString("~gPermission level for ~W%s ~ghas been changed to ~W%s.\r\n"), help.Name, imcperm_names[permvalue])
		help.Level = permvalue
		imc_savehelps()
		return
	}
	imchedit(ch, libc.CString(""))
}
func imcrefresh(ch *char_data, argument *byte) {
	var (
		r     *REMOTEINFO
		rnext *REMOTEINFO
	)
	for r = first_rinfo; r != nil; r = rnext {
		rnext = r.Next
		imc_delete_reminfo(r)
	}
	imc_request_keepalive()
	imc_to_char(libc.CString("Mud list is being refreshed.\r\n"), ch)
}
func imctemplates(ch *char_data, argument *byte) {
	imc_to_char(libc.CString("Refreshing all templates.\r\n"), ch)
	imc_load_templates()
}
func imclast(ch *char_data, argument *byte) {
	var p *IMC_PACKET
	p = imc_newpacket(GET_NAME(ch), libc.CString("imc-laston"), this_imcmud.Servername)
	if argument != nil && *argument != '\x00' {
		imc_addtopacket(p, libc.CString("username=%s"), argument)
	}
	imc_write_packet(p)
}
func imc_other(ch *char_data, argument *byte) {
	var (
		buf  [4096]byte
		cmd  *IMC_CMD_DATA
		col  int
		perm int
	)
	strlcpy(&buf[0], libc.CString("~gThe following commands are available:\r\n"), LGST)
	imcstrlcat(&buf[0], libc.CString("~G-------------------------------------\r\n\r\n"), LGST)
	for perm = IMCPERM_MORT; perm <= ch.Player_specials.Imcchardata.Imcperm; perm++ {
		col = 0
		stdio.Snprintf(&buf[C.strlen(&buf[0])], int(LGST-C.strlen(&buf[0])), "\r\n~g%s commands:~G\r\n", imcperm_names[perm])
		for cmd = first_imc_command; cmd != nil; cmd = cmd.Next {
			if cmd.Level != perm {
				continue
			}
			stdio.Snprintf(&buf[C.strlen(&buf[0])], int(LGST-C.strlen(&buf[0])), "%-15s", cmd.Name)
			if func() int {
				p := &col
				*p++
				return *p
			}()%6 == 0 {
				imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
			}
		}
		if col%6 != 0 {
			imcstrlcat(&buf[0], libc.CString("\r\n"), LGST)
		}
	}
	imc_to_pager(&buf[0], ch)
	imc_to_pager(libc.CString("\r\n~gFor information about a specific command, see ~Wimchelp <command>~g.\r\n"), ch)
}
func imc_find_social(ch *char_data, sname *byte, person *byte, mud *byte, victim int) *byte {
	var (
		socname [4096]byte
		social  *social_messg
		i       int
	)
	socname[0] = '\x00'
	var lcSocName [4096]byte
	for i = 0; i < LGST; i++ {
		if *(*byte)(unsafe.Add(unsafe.Pointer(sname), i)) != '\x00' {
			lcSocName[i] = byte(int8(C.tolower(int(*(*byte)(unsafe.Add(unsafe.Pointer(sname), i))))))
		} else {
			lcSocName[i] = '\x00'
		}
	}
	if (func() *social_messg {
		social = find_social(&lcSocName[0])
		return social
	}()) == nil {
		imc_printf(ch, libc.CString("~YSocial ~W%s~Y does not exist on this mud.\r\n"), sname)
		return &socname[0]
	}
	if person != nil && *person != '\x00' && mud != nil && *mud != '\x00' {
		if person != nil && *person != '\x00' && C.strcasecmp(person, GET_NAME(ch)) == 0 && mud != nil && *mud != '\x00' && C.strcasecmp(mud, this_imcmud.Localname) == 0 {
			if social.Others_auto == nil {
				imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing others_auto.\r\n"), sname)
				return &socname[0]
			}
			imcstrlcpy(&socname[0], social.Others_auto, LGST)
		} else {
			if victim == 0 {
				if social.Others_found == nil {
					imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing others_found.\r\n"), sname)
					return &socname[0]
				}
				imcstrlcpy(&socname[0], social.Others_found, LGST)
			} else if victim == 1 {
				if social.Vict_found == nil {
					imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing vict_found.\r\n"), sname)
					return &socname[0]
				}
				imcstrlcpy(&socname[0], social.Vict_found, LGST)
			} else {
				if social.Char_found == nil {
					imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing char_found.\r\n"), sname)
					return &socname[0]
				}
				imcstrlcpy(&socname[0], social.Char_found, LGST)
			}
		}
	} else {
		if victim == 0 || victim == 1 {
			if social.Others_no_arg == nil {
				imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing others_no_arg.\r\n"), sname)
				return &socname[0]
			}
			imcstrlcpy(&socname[0], social.Others_no_arg, LGST)
		} else {
			if social.Char_no_arg == nil {
				imc_printf(ch, libc.CString("~YSocial ~W%s~Y: Missing char_no_arg.\r\n"), sname)
				return &socname[0]
			}
			imcstrlcpy(&socname[0], social.Char_no_arg, LGST)
		}
	}
	return &socname[0]
}
func imc_act_string(format *byte, ch *char_data, vic *char_data) *byte {
	var (
		he_she       [3]*byte = [3]*byte{libc.CString("it"), libc.CString("he"), libc.CString("she")}
		him_her      [3]*byte = [3]*byte{libc.CString("it"), libc.CString("him"), libc.CString("her")}
		his_her      [3]*byte = [3]*byte{libc.CString("its"), libc.CString("his"), libc.CString("her")}
		buf          [4096]byte
		tmp_str      [4096]byte
		i            *byte = libc.CString("")
		point        *byte
		should_upper bool = FALSE != 0
	)
	if format == nil || *format == '\x00' || ch == nil {
		return nil
	}
	point = &buf[0]
	for *format != '\x00' {
		if *format == '.' || *format == '?' || *format == '!' {
			should_upper = TRUE != 0
		} else if int(libc.BoolToInt(should_upper)) == TRUE && (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*format)))))&int(uint16(int16(_ISspace)))) == 0 && *format != '$' {
			should_upper = FALSE != 0
		}
		if *format != '$' {
			*func() *byte {
				p := &point
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}() = *func() *byte {
				p := &format
				x := *p
				*p = (*byte)(unsafe.Add(unsafe.Pointer(*p), 1))
				return x
			}()
			continue
		}
		format = (*byte)(unsafe.Add(unsafe.Pointer(format), 1))
		if vic == nil && (*format == 'N' || *format == 'E' || *format == 'M' || *format == 'S' || *format == 'K') {
			i = libc.CString(" !!!!! ")
		} else {
			switch *format {
			default:
				i = libc.CString(" !!!!! ")
			case 'n':
				i = imc_makename(GET_NAME(ch), this_imcmud.Localname)
			case 'N':
				i = GET_NAME(vic)
			case 'e':
				if should_upper {
					i = imccapitalize(he_she[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()])
				} else {
					i = he_she[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()]
				}
			case 'E':
				if should_upper {
					i = imccapitalize(he_she[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()])
				} else {
					i = he_she[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()]
				}
			case 'm':
				if should_upper {
					i = imccapitalize(him_her[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()])
				} else {
					i = him_her[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()]
				}
			case 'M':
				if should_upper {
					i = imccapitalize(him_her[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()])
				} else {
					i = him_her[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()]
				}
			case 's':
				if should_upper {
					i = imccapitalize(his_her[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()])
				} else {
					i = his_her[func() int8 {
						if ch.Sex < 0 {
							return 0
						}
						if ch.Sex > 2 {
							return 2
						}
						return ch.Sex
					}()]
				}
			case 'S':
				if should_upper {
					i = imccapitalize(his_her[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()])
				} else {
					i = his_her[func() int8 {
						if vic.Sex < 0 {
							return 0
						}
						if vic.Sex > 2 {
							return 2
						}
						return vic.Sex
					}()]
				}
			case 'k':
				imcone_argument(GET_NAME(ch), &tmp_str[0])
				i = &tmp_str[0]
			case 'K':
				imcone_argument(GET_NAME(vic), &tmp_str[0])
				i = &tmp_str[0]
			}
		}
		format = (*byte)(unsafe.Add(unsafe.Pointer(format), 1))
		for (func() byte {
			p := point
			*point = *i
			return *p
		}()) != '\x00' {
			point = (*byte)(unsafe.Add(unsafe.Pointer(point), 1))
			i = (*byte)(unsafe.Add(unsafe.Pointer(i), 1))
		}
	}
	*point = 0
	point = (*byte)(unsafe.Add(unsafe.Pointer(point), 1))
	*point = '\x00'
	buf[0] = byte(int8(C.toupper(int(buf[0]))))
	return &buf[0]
}
func imc_make_skeleton(name *byte) *char_data {
	var skeleton *char_data
	for {
		if (func() *char_data {
			skeleton = new(char_data)
			return skeleton
		}()) == nil {
			imclog(libc.CString("Malloc failure @ %s:%d\n"), __FILE__, __LINE__)
			abort()
		}
		if true {
			break
		}
	}
	skeleton.Name = C.strdup(name)
	skeleton.Short_descr = C.strdup(name)
	skeleton.In_room = real_room(1)
	return skeleton
}
func imc_purge_skeleton(skeleton *char_data) {
	if skeleton == nil {
		return
	}
	for {
		if skeleton.Name != nil {
			libc.Free(unsafe.Pointer(skeleton.Name))
			skeleton.Name = nil
		}
		if true {
			break
		}
	}
	for {
		if skeleton.Short_descr != nil {
			libc.Free(unsafe.Pointer(skeleton.Short_descr))
			skeleton.Short_descr = nil
		}
		if true {
			break
		}
	}
	for {
		if skeleton != nil {
			libc.Free(unsafe.Pointer(skeleton))
			skeleton = nil
		}
		if true {
			break
		}
	}
}
func imc_send_social(ch *char_data, argument *byte, telloption int) *byte {
	var (
		skeleton *char_data = nil
		ps       *byte
		msg      [4096]byte
		socbuf   [4096]byte
		arg1     [1024]byte
		person   [1024]byte
		mud      [1024]byte
		buf      [4096]byte
		x        uint64
	)
	person[0] = '\x00'
	mud[0] = '\x00'
	argument = imcone_argument(argument, &arg1[0])
	if argument != nil && *argument != '\x00' {
		if (func() *byte {
			ps = C.strchr(argument, '@')
			return ps
		}()) == nil {
			imc_to_char(libc.CString("You need to specify a person@mud for a target.\r\n"), ch)
			return libc.CString("")
		} else {
			for x = 0; x < uint64(C.strlen(argument)); x++ {
				person[x] = *(*byte)(unsafe.Add(unsafe.Pointer(argument), x))
				if person[x] == '@' {
					break
				}
			}
			person[x] = '\x00'
			*ps = '\x00'
			strlcpy(&mud[0], (*byte)(unsafe.Add(unsafe.Pointer(ps), 1)), SMST)
		}
	}
	if telloption == 0 {
		stdio.Snprintf(&socbuf[0], LGST, "%s", imc_find_social(ch, &arg1[0], &person[0], &mud[0], 0))
		if socbuf[0] == '\x00' {
			return libc.CString("")
		}
	}
	if telloption == 1 {
		stdio.Snprintf(&socbuf[0], LGST, "%s", imc_find_social(ch, &arg1[0], &person[0], &mud[0], 1))
		if socbuf[0] == '\x00' {
			return libc.CString("")
		}
	}
	if telloption == 2 {
		stdio.Snprintf(&socbuf[0], LGST, "%s", imc_find_social(ch, &arg1[0], &person[0], &mud[0], 2))
		if socbuf[0] == '\x00' {
			return libc.CString("")
		}
	}
	if argument != nil && *argument != '\x00' {
		var sex int
		stdio.Snprintf(&buf[0], LGST, "%s@%s", &person[0], &mud[0])
		sex = imc_get_ucache_gender(&buf[0])
		if sex == -1 {
			imc_send_ucache_request(&buf[0])
			sex = SEX_MALE
		} else {
			sex = imctodikugender(sex)
		}
		skeleton = imc_make_skeleton(&buf[0])
		skeleton.Sex = int8(sex)
	}
	stdio.Snprintf(&msg[0], LGST, "%s", imc_act_string(color_mtoi(&socbuf[0]), ch, skeleton))
	if skeleton != nil {
		imc_purge_skeleton(skeleton)
	}
	return &msg[0]
}
func imc_funcname(func_ IMC_FUN) *byte {
	if libc.FuncAddr(func_) == libc.FuncAddr(imc_other) {
		return libc.CString("imc_other")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imclisten) {
		return libc.CString("imclisten")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcchanlist) {
		return libc.CString("imcchanlist")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imclist) {
		return libc.CString("imclist")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcinvis) {
		return libc.CString("imcinvis")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcwho) {
		return libc.CString("imcwho")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imclocate) {
		return libc.CString("imclocate")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imctell) {
		return libc.CString("imctell")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcreply) {
		return libc.CString("imcreply")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcbeep) {
		return libc.CString("imcbeep")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcignore) {
		return libc.CString("imcignore")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcfinger) {
		return libc.CString("imcfinger")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcinfo) {
		return libc.CString("imcinfo")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imccolor) {
		return libc.CString("imccolor")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcafk) {
		return libc.CString("imcafk")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcchanwho) {
		return libc.CString("imcchanwho")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcconnect) {
		return libc.CString("imcconnect")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcdisconnect) {
		return libc.CString("imcdisconnect")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcpermstats) {
		return libc.CString("imcpermstats")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imc_deny_channel) {
		return libc.CString("imc_deny_channel")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcpermset) {
		return libc.CString("imcpermset")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcsetup) {
		return libc.CString("imcsetup")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imccommand) {
		return libc.CString("imccommand")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcban) {
		return libc.CString("imcban")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcconfig) {
		return libc.CString("imcconfig")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imc_show_ucache_contents) {
		return libc.CString("imc_show_ucache_contents")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcremoteadmin) {
		return libc.CString("imcremoteadmin")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcdebug) {
		return libc.CString("imcdebug")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imchedit) {
		return libc.CString("imchedit")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imchelp) {
		return libc.CString("imchelp")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imccedit) {
		return libc.CString("imccedit")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imcrefresh) {
		return libc.CString("imcrefresh")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imclast) {
		return libc.CString("imclast")
	}
	if libc.FuncAddr(func_) == libc.FuncAddr(imctemplates) {
		return libc.CString("imctemplates")
	}
	return libc.CString("")
}
func imc_function(func_ *byte) IMC_FUN {
	if C.strcasecmp(func_, libc.CString("imc_other")) == 0 {
		return func(ch *char_data, argument *byte) {
			imc_other(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imclisten")) == 0 {
		return func(ch *char_data, argument *byte) {
			imclisten(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcchanlist")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcchanlist(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imclist")) == 0 {
		return func(ch *char_data, argument *byte) {
			imclist(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcinvis")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcinvis(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcwho")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcwho(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imclocate")) == 0 {
		return func(ch *char_data, argument *byte) {
			imclocate(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imctell")) == 0 {
		return func(ch *char_data, argument *byte) {
			imctell(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcreply")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcreply(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcbeep")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcbeep(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcignore")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcignore(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcfinger")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcfinger(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcinfo")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcinfo(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imccolor")) == 0 {
		return func(ch *char_data, argument *byte) {
			imccolor(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcafk")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcafk(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcchanwho")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcchanwho(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcconnect")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcconnect(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcdisconnect")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcdisconnect(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcpermstats")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcpermstats(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imc_deny_channel")) == 0 {
		return func(ch *char_data, argument *byte) {
			imc_deny_channel(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcpermset")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcpermset(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcsetup")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcsetup(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imccommand")) == 0 {
		return func(ch *char_data, argument *byte) {
			imccommand(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcban")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcban(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcconfig")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcconfig(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imc_show_ucache_contents")) == 0 {
		return func(ch *char_data, argument *byte) {
			imc_show_ucache_contents(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcremoteadmin")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcremoteadmin(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcdebug")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcdebug(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imchelp")) == 0 {
		return func(ch *char_data, argument *byte) {
			imchelp(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imccedit")) == 0 {
		return func(ch *char_data, argument *byte) {
			imccedit(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imchedit")) == 0 {
		return func(ch *char_data, argument *byte) {
			imchedit(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imcrefresh")) == 0 {
		return func(ch *char_data, argument *byte) {
			imcrefresh(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imclast")) == 0 {
		return func(ch *char_data, argument *byte) {
			imclast(ch, argument)
		}
	}
	if C.strcasecmp(func_, libc.CString("imctemplates")) == 0 {
		return func(ch *char_data, argument *byte) {
			imctemplates(ch, argument)
		}
	}
	return nil
}
func imc_command_hook(ch *char_data, command *byte, argument *byte) bool {
	var (
		cmd   *IMC_CMD_DATA
		alias *IMC_ALIAS
		c     *IMC_CHANNEL
		p     *byte
		arg   *byte
	)
	skip_spaces(&argument)
	if IS_NPC(ch) {
		return FALSE != 0
	}
	if this_imcmud == nil {
		imcbug(libc.CString("%s"), "Ooops. IMC being called with no configuration!")
		return FALSE != 0
	}
	if first_imc_command == nil {
		imcbug(libc.CString("%s"), "ACK! There's no damn command data loaded!")
		return FALSE != 0
	}
	if ch.Player_specials.Imcchardata.Imcperm <= IMCPERM_NONE {
		return FALSE != 0
	}
	arg = argument
	skip_spaces(&arg)
	for cmd = first_imc_command; cmd != nil; cmd = cmd.Next {
		if ch.Player_specials.Imcchardata.Imcperm < cmd.Level {
			continue
		}
		for alias = cmd.First_alias; alias != nil; alias = alias.Next {
			if C.strcasecmp(command, alias.Name) == 0 {
				command = cmd.Name
				break
			}
		}
		if C.strcasecmp(command, cmd.Name) == 0 {
			if int(libc.BoolToInt(cmd.Connected)) == TRUE && int(this_imcmud.State) < IMC_ONLINE {
				imc_to_char(libc.CString("The mud is not currently connected to IMC2.\r\n"), ch)
				return TRUE != 0
			}
			if cmd.Function == nil {
				imc_to_char(libc.CString("That command has no code set. Inform the administration.\r\n"), ch)
				imcbug(libc.CString("imc_command_hook: Command %s has no code set!"), cmd.Name)
				return TRUE != 0
			}
			(cmd.Function)(ch, argument)
			return TRUE != 0
		}
	}
	c = imc_findchannel(command)
	if c == nil || int(c.Level) > ch.Player_specials.Imcchardata.Imcperm {
		return FALSE != 0
	}
	if imc_hasname(ch.Player_specials.Imcchardata.Imc_denied, c.Local_name) {
		imc_printf(ch, libc.CString("You have been denied the use of %s by the administration.\r\n"), c.Local_name)
		return TRUE != 0
	}
	if !c.Refreshed {
		imc_printf(ch, libc.CString("The %s channel has not yet been refreshed by the server.\r\n"), c.Local_name)
		return TRUE != 0
	}
	if argument == nil || *argument == '\x00' {
		var y int
		imc_printf(ch, libc.CString("~cThe last %d %s messages:\r\n"), MAX_IMCHISTORY, c.Local_name)
		for y = 0; y < MAX_IMCHISTORY; y++ {
			if c.History[y] != nil {
				imc_printf(ch, libc.CString("%s\r\n"), c.History[y])
			} else {
				break
			}
		}
		return TRUE != 0
	}
	if ch.Player_specials.Imcchardata.Imcperm >= IMCPERM_ADMIN && C.strcasecmp(argument, libc.CString("log")) == 0 {
		if (c.Flags & (1 << 0)) == 0 {
			c.Flags |= 1 << 0
			imc_printf(ch, libc.CString("~RFile logging enabled for %s, PLEASE don't forget to undo this when it isn't needed!\r\n"), c.Local_name)
		} else {
			c.Flags &= ^(1 << 0)
			imc_printf(ch, libc.CString("~GFile logging disabled for %s.\r\n"), c.Local_name)
		}
		imc_save_channels()
		return TRUE != 0
	}
	if !imc_hasname(ch.Player_specials.Imcchardata.Imc_listen, c.Local_name) {
		imc_printf(ch, libc.CString("You are not currently listening to %s. Use the imclisten command to listen to this channel.\r\n"), c.Local_name)
		return TRUE != 0
	}
	switch *argument {
	case ',':
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		imc_sendmessage(c, GET_NAME(ch), color_mtoi(argument), 1)
	case '@':
		argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		for (int(*(*uint16)(unsafe.Add(unsafe.Pointer(*__ctype_b_loc()), unsafe.Sizeof(uint16(0))*uintptr(int(*argument))))) & int(uint16(int16(_ISspace)))) != 0 {
			argument = (*byte)(unsafe.Add(unsafe.Pointer(argument), 1))
		}
		p = imc_send_social(ch, argument, 0)
		if p == nil || *p == '\x00' {
			return TRUE != 0
		}
		imc_sendmessage(c, GET_NAME(ch), p, 2)
	default:
		imc_sendmessage(c, GET_NAME(ch), color_mtoi(argument), 0)
	}
	return TRUE != 0
}
