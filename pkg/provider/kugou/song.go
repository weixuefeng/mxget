package kugou

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"

	"github.com/winterssy/ghttp"
	"github.com/winterssy/mxget/pkg/api"
)

func (a *API) GetSong(ctx context.Context, hash string) (*api.SongResponse, error) {
	resp, err := a.GetSongRaw(ctx, hash)
	if err != nil {
		return nil, err
	}

	a.patchSongsInfo(ctx, &resp.Song)
	a.patchSongsLyric(ctx, &resp.Song)
	songs := resolve(&resp.Song)
	return songs[0], nil
}

// 获取歌曲详情
func (a *API) GetSongRaw(ctx context.Context, hash string) (*SongResponse, error) {
	params := ghttp.Params{
		"hash": hash,
	}

	resp := new(SongResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSong,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.ErrCode != 0 {
		return nil, fmt.Errorf("get song: %s", resp.errorMessage())
	}

	return resp, nil
}

func (a *API) GetSongURL(ctx context.Context, hash string) (string, error) {
	resp, err := a.GetSongURLRaw(ctx, hash)
	if err != nil {
		return "", err
	}
	if len(resp.URL) == 0 {
		return "", errors.New("get song url: no data")
	}

	return resp.URL[rand.Intn(len(resp.URL))], nil
}

// 获取歌曲播放地址
func (a *API) GetSongURLRaw(ctx context.Context, hash string) (*SongURLResponse, error) {
	data := []byte(hash + "kgcloudv2")
	key := fmt.Sprintf("%x", md5.Sum(data))

	params := ghttp.Params{
		"hash": hash,
		"key":  key,
	}

	resp := new(SongURLResponse)
	err := a.SendRequest(ghttp.MethodGet, APIGetSongURL,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != 1 {
		if resp.Status == 2 {
			err = errors.New("get song url: copyright protection")
		} else {
			err = fmt.Errorf("get song url: %d", resp.Status)
		}
		return nil, err
	}

	return resp, nil
}

// 获取歌词
func (a *API) GetSongLyric(ctx context.Context, hash string) (string, error) {
	params := ghttp.Params{
		"hash": hash,
	}
	return a.SendRequest(ghttp.MethodGet, APIGetSongLyric,
		ghttp.WithQuery(params),
		ghttp.WithContext(ctx),
	).Text()
}
