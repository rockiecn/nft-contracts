// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";

contract DmNFT is ERC721Enumerable, Ownable {
    string private _customBaseURI;
    uint256 public nextTokenId = 1;

    // 构造函数
    constructor(string memory name_, string memory symbol_) 
        ERC721(name_, symbol_) 
        Ownable(msg.sender)  // 明确设置合约部署者为owner
    {}

    function _baseURI() internal view override returns (string memory) {
        return _customBaseURI;
    }

    function setBaseURI(string memory baseURI_) public onlyOwner {
        _customBaseURI = baseURI_;
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId); // 替代 _exists 检查
        return string(abi.encodePacked(_baseURI(), Strings.toString(tokenId)));
    }

    function mint(address to) public {
        _safeMint(to, nextTokenId);
        nextTokenId++;
    }

    // 获取用户拥有的所有 tokenId
    function tokensOfOwner(address owner) public view returns (uint256[] memory) {
        uint256 balance = balanceOf(owner);
        uint256[] memory tokens = new uint256[](balance);
        
        for (uint256 i = 0; i < balance; i++) {
            tokens[i] = tokenOfOwnerByIndex(owner, i);
        }
        
        return tokens;
    }
}