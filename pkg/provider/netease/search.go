package netease

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) SearchSongs(ctx context.Context, keyword string) (*api.SearchSongsResponse, error) {
	resp, err := a.SearchSongsRaw(ctx, keyword, 0, 50)
	if err != nil {
		return nil, err
	}

	n := len(resp.Result.Songs)
	if n == 0 {
		return nil, errors.New("search songs: no data")
	}

	songs := make([]*api.SongResponse, n)
	for i, s := range resp.Result.Songs {
		artists := make([]string, len(s.Artists))
		for j, a := range s.Artists {
			artists[j] = strings.TrimSpace(a.Name)
		}
		songs[i] = &api.SongResponse{
			Id:     strconv.Itoa(s.Id),
			Name:   strings.TrimSpace(s.Name),
			Artist: strings.Join(artists, "/"),
			Album:  s.Album.Name,
		}
	}
	return &api.SearchSongsResponse{
		Keyword: keyword,
		Count:   uint32(n),
		Songs:   songs,
	}, nil
}

// 搜索歌曲
func (a *API) SearchSongsRaw(ctx context.Context, keyword string, offset int, limit int) (*SearchSongsResponse, error) {
	// type: 1: 单曲, 10: 专辑, 100: 歌手, 1000: 歌单, 1002: 用户,
	// 1004: MV, 1006: 歌词, 1009: 电台, 1014: 视频
	data := map[string]interface{}{
		"s":      keyword,
		"type":   1,
		"offset": offset,
		"limit":  limit,
	}

	resp := new(SearchSongsResponse)
	err := a.SendRequest(ghttp.MethodPost, APISearch,
		ghttp.WithForm(weapi(data)),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Code != 200 {
		return nil, fmt.Errorf("search songs: %s", resp.errorMessage())
	}

	return resp, nil
}
